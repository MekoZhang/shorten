package short

import (
	"errors"
	"log"

	"github.com/zhangxd1989/shorten/base"
	"github.com/zhangxd1989/shorten/conf"
	"github.com/zhangxd1989/shorten/sequence"
	"github.com/zhangxd1989/shorten/redis"
	_ "github.com/zhangxd1989/shorten/sequence/db"
)

type shorter struct {
	sequence sequence.Sequence
}

// initSequence will panic when it can not open the sequence successfully.
func (shorter *shorter) initSequence() {
	seq, err := sequence.GetSequence("db")
	if err != nil {
		log.Panicf("get sequence instance error. %v", err)
	}

	shorter.sequence = seq
}

func (shorter *shorter) Expand(shortURL string) (longURL string, err error) {

	longURL, err = redis.Get(shortURL)
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

	err = redis.Set(shortURL, longURL)
	if err != nil {
		log.Printf("short db set error. %v", err)
		return "", errors.New("short db set error")
	}

	return shortURL, nil
}

var Shorter shorter

func Start() {
	Shorter.initSequence()
	log.Println("shorter starts")
}
