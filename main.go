package main

import (
	"github.com/ChangSZ/golib/log"

	"github.com/ChangSZ/blog/conf"
	"github.com/ChangSZ/blog/router"
)

func main() {
	conf.DefaultInit()

	r := router.RoutersInit()

	log.Info("app Run...")
	if err := r.Run(":8081"); err != nil {
		panic(err)
	}
}
