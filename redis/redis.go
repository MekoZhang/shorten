package redis

import (
	"fmt"

	"github.com/zhangxd1989/shorten/conf"
	"github.com/garyburd/redigo/redis"
	"time"
	"os"
	"os/signal"
	"syscall"
)

var (
	Pool *redis.Pool
)

func init() {
	Pool = newPool(conf.Conf.Redis.Port)
	close()
}

func newPool(server string) *redis.Pool {

	return &redis.Pool{

		MaxIdle:     conf.Conf.Redis.MaxIdle,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func close() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		Pool.Close()
		os.Exit(0)
	}()
}

func get(key string) (string, error) {

	conn := Pool.Get()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	dataString := string(data)
	if err != nil {
		return dataString, fmt.Errorf("error get key %s: %v", key, err)
	}
	return dataString, err
}

func set(key, value string) (error) {

	conn := Pool.Get()
	defer conn.Close()

	_, err := redis.Bytes(conn.Do("SET", key, value))
	if err != nil {
		fmt.Println("redis set failed:", err)
	}
	if err != nil {
		return fmt.Errorf("error get key %s: %v", key, err)
	}
	return err
}
