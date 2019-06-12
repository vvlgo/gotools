package redisclient

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"io"
	"strings"
	"time"
)

func newPool(addr string, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     100,
		IdleTimeout: 240 * time.Second,
		// Other pool configuration not shown in this example.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, nil
		},
	}
}

var (
	redisPool *redis.Pool
)

/*
RedisInit redis 初始化
*/
func RedisInit(host, port, password string) {
	redisPool = newPool(host+":"+port, password)
}

type MyRedisReConn struct {
	RedisConn redis.Conn
	ReidsDB   int
}

/*
GetRedisConn 从redis池子拿去一个n号库的链接
n int redis 对应的n号库
*/
func GetRedisConn(n int) (conn *MyRedisReConn, err error) {
	conn = &MyRedisReConn{}
	conn.RedisConn = redisPool.Get()
	_, err = conn.RedisConn.Do("SELECT", n)
	if err != nil {
		return nil, err
	}
	conn.ReidsDB = n
	return conn, nil
}

/*
IsConnError 判断是否有错误
*/
func IsConnError(err error) bool {
	var needNewConn bool

	if err == nil {
		return false
	}

	if err == io.EOF {
		needNewConn = true
	}
	if strings.Contains(err.Error(), "use of closed network connection") {
		needNewConn = true
	}
	if strings.Contains(err.Error(), "connect: connection refused") {
		needNewConn = true
	}
	if strings.Contains(err.Error(), "connection closed") {
		needNewConn = true
	}
	return needNewConn
}

/*
Redo 在pool加入TestOnBorrow方法来去除扫描坏连接，并重新连接redis
*/
func (myredis MyRedisReConn) Redo(command string, opt ...interface{}) (interface{}, error) {
	defer myredis.RedisConn.Close()
	var conn redis.Conn
	var err error
	var maxretry = 3
	var needNewConn bool
	resp, err := myredis.RedisConn.Do(command, opt...)
	needNewConn = IsConnError(err)
	if !needNewConn {
		return resp, err
	} else {
		conn, err = redisPool.Dial()
		_, err = conn.Do("SELECT", myredis.ReidsDB)
	}
	for index := 0; index < maxretry; index++ {
		if conn == nil && index+1 > maxretry {
			return resp, err
		}
		if conn == nil {
			conn, err = redisPool.Dial()
			_, err = conn.Do("SELECT", myredis.ReidsDB)
		}
		if err != nil {
			continue
		}
		resp, err := conn.Do(command, opt...)
		needNewConn = IsConnError(err)
		if !needNewConn {
			return resp, err
		} else {
			conn, err = redisPool.Dial()
			_, err = conn.Do("SELECT", myredis.ReidsDB)
		}
	}
	conn.Close()
	return "", errors.New("redis error")
}
