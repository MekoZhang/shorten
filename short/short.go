package short

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/zhangxd1989/shorten/base"
	"github.com/zhangxd1989/shorten/conf"
	"github.com/zhangxd1989/shorten/sequence"
	"github.com/garyburd/redigo/redis"
	_ "github.com/zhangxd1989/shorten/sequence/db"
	"time"
	"os"
	"os/signal"
	"syscall"
)

type shorter struct {
	readDB   *sql.DB
	writeDB  *sql.DB
	sequence sequence.Sequence
}

var (
	Pool *redis.Pool
)

func (shorter *shorter) connect() {
	redisHost := ":6379"
	Pool = newPool(redisHost)
	close()
}

func newPool(server string) *redis.Pool {

	return &redis.Pool{

		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

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

// initSequence will panic when it can not open the sequence successfully.
func (shorter *shorter) initSequence() {
	seq, err := sequence.GetSequence("db")
	if err != nil {
		log.Panicf("get sequence instance error. %v", err)
	}

	err = seq.Open()
	if err != nil {
		log.Panicf("open sequence instance error. %v", err)
	}

	shorter.sequence = seq
}

func (shorter *shorter) Expand(shortURL string) (longURL string, err error) {

	longURL, err = get(shortURL)
	if err != nil {
		log.Printf("short db get error. %v", err)
		return "", errors.New("short db get error")
	}

	return longURL, nil
}

func (shorter *shorter) Short(longURL string) (shortURL string, err error) {
	for {
		var seq uint64
		seq, err = shorter.sequence.NextSequence()
		if err != nil {
			log.Printf("get next sequence error. %v", err)
			return "", errors.New("get next sequence error")
		}

		shortURL = base.Int2String(seq)
		if _, exists := conf.Conf.Common.BlackShortURLsMap[shortURL]; exists {
			continue
		} else {
			break
		}
	}

	err = set(shortURL, longURL)
	if err != nil {
		log.Printf("short db set error. %v", err)
		return "", errors.New("short db set error")
	}

	return shortURL, nil
}

var Shorter shorter

func Start() {
	Shorter.connect()
	Shorter.initSequence()
	log.Println("shorter starts")
}
