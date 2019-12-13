package redis

import (
	"encoding"
	"strings"
	"time"
)

type KeyDetail struct {
	Key   string        `json:"key"`
	Type  string        `json:"type"`
	TTL   time.Duration `json:"ttl"`
	Value interface{}   `json:"value"`
}

type Collection struct {
	Total int64    `json:"total"`
	Keys  []string `json:"keys"`
}

type Info map[string]interface{}

var _ encoding.BinaryUnmarshaler = (*Info)(nil)

func (r Info) UnmarshalBinary(data []byte) error {
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		keyVal := strings.Split(line, ":")
		key := strings.TrimSpace(keyVal[0])
		val := strings.TrimSpace(keyVal[1])
		r[key] = val
	}
	return nil
}

func (r Info) GetString(key string) string {
	v := r[key]
	s, _ := v.(string)
	return s
}

type Menu struct {
	Label  string `json:"label"`
	IsLeaf bool   `json:"isLeaf"`
}

type byMenu []Menu

func (b byMenu) Len() int {
	return len(b)
}

func (b byMenu) Less(i, j int) bool {
	return b[i].Label < b[j].Label
}

func (b byMenu) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
