package db

import (
	"log"

	"github.com/zhangxd1989/shorten/sequence"
	"github.com/zhangxd1989/shorten/redis"
)

type SequenceDB struct {
}

func (dbSeq *SequenceDB) NextSequence() (sequence uint64, err error) {

	// 兼容LastInsertId方法的返回值
	var lastID int64
	lastID, err = redis.Incr("sequence")
	if err != nil {
		log.Printf("sequence db get LastInsertId error. %v", err)
		return 0, err
	} else {
		sequence = uint64(lastID)
		// started at 0. :)
		sequence -= 1
		return sequence, nil
	}
}

var dbSeq = SequenceDB{}

func init() {
	sequence.Register("db", &dbSeq)
}
