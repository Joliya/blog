package main

import (
	"blog/conf"
	"blog/internal/dao"
	"blog/internal/es"
	"blog/internal/routers"
	"blog/pkg/logger"
	"blog/pkg/redis"
	"blog/pkg/shutdown"
	"blog/pkg/utils"
	"context"
	"github.com/spf13/pflag"
	"log"
	"net/http"
	"time"
)

var (
	cfgFile = pflag.StringP("config", "c", "./conf/dev.yml", "config file path.")
	//version = pflag.BoolP("version", "v", false, "show version info.")
)

func main() {
	pflag.Parse()
	cfg := conf.Init(*cfgFile)
	logger.Init(&cfg.Logger)
	redis.Init(&cfg.Redis)
	dao.Init(&cfg.Mysql)
	if !cfg.Elasticsearch.Disable {
		es.Init(&cfg.Elasticsearch)
	}

	addr := cfg.App.Addr
	srv := &http.Server{
		Addr:    addr,
		Handler: routers.InitRouter(),
	}

	go func() {
		if err := srv.ListenAndServe(); utils.IsNotNil(err) {
			log.Println("server run: ", err)
		}
	}()

	shutdown.NewHook().Close(
		// 关闭 http server
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				log.Println("http server closed err", err)
			} else {
				log.Println("http server closed")
			}
		},
	)
}
