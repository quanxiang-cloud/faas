package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/quanxiang-cloud/cabin/logger"
	restful "github.com/quanxiang-cloud/faas/api"
	"github.com/quanxiang-cloud/faas/pkg/config"
)

var (
	configPath = flag.String("config", "configs/config.yml", "-config 配置文件地址")
)

func main() {
	flag.Parse()
	log := logger.Logger
	conf, err := config.NewConfig(*configPath)
	if err != nil {
		panic(err)
	}
	// start
	ctx := context.Background()
	router, err := restful.NewRouter(ctx, conf, log)
	if err != nil {
		panic(err)
	}
	go router.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			router.Close()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
