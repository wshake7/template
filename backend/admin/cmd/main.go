package main

import (
	"admin/config"
	"admin/fiberc"
	"admin/router"
	"admin/services"
	"flag"
	"go-common/log"
	"go-common/viperc"
)

var configFile = flag.String("f", "./etc/config.yaml", "the config file")

func main() {
	flag.Parse()
	conf := config.Conf
	_, err := viperc.ParseFile(*configFile, conf)
	if err != nil {
		panic(err)
	}
	log.Init(log.DevConfig())
	services.New(conf)
	app := fiberc.NewFiber(conf)
	group := app.Group(conf.RestPrefix)
	r := router.Router{}
	r.RegisterRouters(group)
	app.Start()
}
