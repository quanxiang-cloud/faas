package es

import (
	"context"
	"fmt"
	"testing"
	"time"

	"git.internal.yunify.com/qxp/misc/elastic2"
	"git.internal.yunify.com/qxp/misc/logger"
)

func Test_builderLogRepo_Search(t *testing.T) {
	elasticClient, err := elastic2.NewClient(&elastic2.Config{
		Host: []string{"http://192.168.200.45:9200/"},
		Log:  true,
	}, logger.Logger)
	if err != nil {
		return
	}
	blr := NewBuildLogRepo(elasticClient)
	fmt.Printf("time.Now().UTC(): %v\n", time.Now().UTC())

	d, err2 := time.ParseDuration("2021-11-17 07:02:00.978984424 +0000 UTC")
	fmt.Println(d, err2)

	tt := time.Unix(1637132520, 0)
	// tt := time.Unix(0, 0)
	blr.Search(context.Background(), "64410a8f-24c8-4a50-a501-ccc951364760", tt, 2, 10)
}
