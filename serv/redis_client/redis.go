package redis_client

import (
	"context"
	"fmt"
	"github.com/ikaven1024/redisAdmin/util"
	"log"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
)

type Client struct {
	redis.Cmdable
}

func NewRedisCli(addr string, password string, db int) *Client {
	cli := redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: db})
	cli.AddHook(&LogHook{})

	return &Client{
		Cmdable: cli,
	}
}

func NewRedisClusterCli(addrs []string, password string) *Client {
	return &Client{
		Cmdable: redis.NewClusterClient(&redis.ClusterOptions{Addrs: addrs, Password: password}),
	}
}

func (r Client) KeyMenus(prefix string) ([]Menu, error) {
	var keySet = make(map[Menu]struct{})
	var cursor uint64 = 0
	var err error
	var keys []string

	if len(prefix) > 0 && !strings.HasPrefix(prefix, ":") {
		prefix = prefix + ":"
	}
	for {
		keys, cursor, err = r.Scan(cursor, prefix+"*", 10000).Result()
		if err != nil {
			return nil, err
		}

		for _, k := range keys {
			suffix := strings.TrimPrefix(k, prefix)
			splits := strings.SplitN(suffix, ":", 2)

			keySet[Menu{
				Label:  splits[0],
				IsLeaf: len(splits) == 1,
			}] = struct{}{}
		}

		if cursor == 0 {
			break
		}
	}

	result := make([]Menu, 0, len(keySet))
	for k := range keySet {
		result = append(result, k)
	}

	sort.Sort(byMenu(result))
	return result, nil
}

type GetValueOpts struct {
	Type     string
	PageNo   int64
	PageSize int64
	Cursor   uint64
	Match    string
}

func (r Client) Info(section ...string) (Info, error) {
	var info = make(Info)
	return info, r.Cmdable.Info(section...).Scan(info)
}

func (r Client) Config(parameter string) (map[string]interface{}, error) {
	items, err := r.Cmdable.ConfigGet(parameter).Result()
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{}, len(items)/2)
	for i := 0; i < len(items); {
		key := items[i]
		i++
		val := items[i]
		i++

		result[util.MustString(key)] = val
	}
	return result, nil
}

func (r Client) GetValue(key string, opts GetValueOpts) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	result["key"] = key
	var err error
	switch strings.ToUpper(opts.Type) {
	case "STRING":
		result["value"], err = r.Get(key).Result()
		return result, err
	case "LIST":
		result["total"] = r.LLen(key).Val()
		result["pageSize"] = opts.PageSize
		result["pageNo"] = opts.PageNo
		result["value"], err = r.LRange(key, (opts.PageNo-1)*opts.PageSize, opts.PageSize*opts.PageNo-1).Result()
		return result, err
	case "HASH":
		result["total"] = r.HLen(key).Val()
		result["value"], result["cursor"], err = r.valueScan(key, opts.Match, opts.Cursor, opts.PageSize, r.HScan)
		return result, err
	case "SET":
		result["total"] = r.SCard(key).Val()
		result["value"], result["cursor"], err = r.valueScan(key, opts.Match, opts.Cursor, opts.PageSize, r.SScan)
		return result, err
	case "ZSET":
		result["total"] = r.ZCard(key).Val()
		result["value"], result["cursor"], err = r.valueScan(key, opts.Match, opts.Cursor, opts.PageSize, r.ZScan)
		return result, err
	default:
		return nil, fmt.Errorf("unknown key type %v", key)
	}
}

func (r Client) valueScan(key, match string, cursor uint64, size int64, scanner func(string, uint64, string, int64) *redis.ScanCmd) ([]string, uint64, error) {

	if len(match) == 0 {
		match = "*"
	}

	data := make([]string, 0, size)
	for {
		var v []string
		var err error
		v, cursor, err = scanner(key, cursor, match, size).Result()
		if err != nil {
			return nil, 0, fmt.Errorf("scanner %v %v error: %v", key, cursor, err)
		}
		data = append(data, v...)
		// scanner to end
		if cursor == 0 {
			break
		}
		// get enough data stop scanner
		if (int64)(len(data)) >= size/2 {
			break
		}
		// data is not enough, continue to scanning
	}

	sort.Strings(data)
	return data, cursor, nil
}

type LogHook struct {
	uid       string
	startTime time.Time
}

func (r *LogHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	r.uid = uid()
	r.startTime = time.Now()

	builder := &strings.Builder{}
	builder.WriteString("[")
	builder.WriteString(r.uid)
	builder.WriteString("] RunCmd:")
	for i := 0; i < len(cmd.Args()); i++ {
		builder.WriteString(" ")
		builder.WriteString(fmt.Sprintf("%v", cmd.Args()[i]))
	}

	log.Println(builder.String())
	return ctx, nil
}

func (r *LogHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	builder := &strings.Builder{}
	builder.WriteString(fmt.Sprintf("[%v] [%v] ", r.uid, time.Now().Sub(r.startTime)))
	if cmd.Err() != nil {
		builder.WriteString(fmt.Sprintf("Response with error: %v.", cmd.Err()))
	} else {
		builder.WriteString(fmt.Sprintf("Responsed."))
	}

	log.Println(builder.String())
	return nil
}

func (r *LogHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (r *LogHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}

var _ redis.Hook = (*LogHook)(nil)

func uid() string {
	pool := "abcdefghigklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	builder := strings.Builder{}
	for i := 0; i < 8; i++ {
		builder.WriteByte(pool[rand.Intn(len(pool))])
	}
	return builder.String()
}
