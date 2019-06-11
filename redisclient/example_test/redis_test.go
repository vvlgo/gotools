package example_test

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/vvlgo/gotools/redisclient"
	"testing"
)

func TestRedis(t *testing.T) {
	redisclient.RedisInit("127.0.0.1", "6379", "")
	reConn, err := redisclient.GetRedisConn(0)
	if err != nil {
		panic(err)
	}
	_, err = reConn.Redo("SET", "key", "value", "EX", 60)
	if err != nil {
		panic(err)
	}
	s, err := redis.String(reConn.Redo("GET", "key"))
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
}
