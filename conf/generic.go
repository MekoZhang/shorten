package conf

import (
	"log"
	"runtime"
)

var Version string

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
