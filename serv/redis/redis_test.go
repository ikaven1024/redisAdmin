package redis

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/go-redis/redis/v7"
)

var cli = NewRedisCli("192.168.11.128:6379")

func TestCreateData(t *testing.T) {
	defer timing(t)()

	num := 1000000
	for i := 0; i < num; i++ {
		cli.Set(fmt.Sprintf("test_key:%v:%v", rand.Intn(1000), rand.Intn(1000)), i, time.Duration(0))
	}

	for i := 0; i < num; i++ {
		cli.HSet("test_large_hash", fmt.Sprintf("%v", i), i)
	}

	for i := 0; i < num; i++ {
		cli.ZAdd("test_large_zset", &redis.Z{Score: float64(i), Member: i})
	}

	for i := 0; i < num; i++ {
		cli.SAdd("test_large_set", i)
	}

	for i := 0; i < num; i++ {
		cli.LPush("test_large_list", i)
	}
}

func TestRedisCli_Keys(t *testing.T) {
	//defer timing(t)()
	//
	//result, err := cli.Keys("*")
	//if err != nil {
	//	t.Error(err)
	//}
	//
	//if len(result) == 0{
	//	t.Error("shall has some result")
	//}
}

func TestRedisCli_Info(t *testing.T) {
	result, err := cli.Info("Server")
	if err != nil {
		t.Error(err)
	}

	version := result.GetString("redis_version")
	if len(version) == 0 {
		t.Error("shall has some result")
	}
}

func TestRedisCli_Config(t *testing.T) {
	result, err := cli.Config("dir")
	if err != nil {
		t.Error(err)
	}

	t.Log(result)
}

func timing(t *testing.T) func() {
	startTime := time.Now()
	return func() {
		t.Logf("takes %v", time.Now().Sub(startTime))
	}
}
