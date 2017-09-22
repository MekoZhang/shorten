package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zhangxd1989/shorten/conf"
	"github.com/zhangxd1989/shorten/redis"
	"github.com/zhangxd1989/shorten/short"
	"github.com/zhangxd1989/shorten/web"
)

func main() {
	cfgFile := flag.String("c", "config.conf", "configuration file")
	version := flag.Bool("v", false, "Version")

	flag.Parse()

	if *version {
		fmt.Println(conf.Version)
		os.Exit(0)
	}

	// parse config
	conf.ParseConfig(*cfgFile)

	// start redis connection
	redis.Start()

	// short service
	short.Start()

	// api
	web.Start()
}
