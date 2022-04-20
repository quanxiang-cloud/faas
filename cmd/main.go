package main

import (
	"context"
	"flag"
	"github.com/quanxiang-cloud/faas/internal/models"
	"github.com/quanxiang-cloud/faas/pkg/k8s"
	"os"
	"os/signal"
	"syscall"

	"github.com/quanxiang-cloud/cabin/logger"
	elastic2 "github.com/quanxiang-cloud/cabin/tailormade/db/elastic"
	mysql2 "github.com/quanxiang-cloud/cabin/tailormade/db/mysql"
	redis2 "github.com/quanxiang-cloud/cabin/tailormade/db/redis"
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
	redisClient, err := redis2.NewClient(conf.Redis)
	if err != nil {
		panic(err)
	}
	esClient, err := elastic2.NewClient(&conf.Elastic, log)
	if err != nil {
		panic(err)
	}
	db, err := mysql2.New(conf.Mysql, log)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(models.Function{})
	k8sClient, err := k8s.NewClient(conf.K8s.NameSpace)
	// start
	ctx := context.Background()
	router, err := restful.NewRouter(ctx, conf, log, db, k8sClient, redisClient, esClient)
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
