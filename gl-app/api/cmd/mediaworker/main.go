package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"gl-app/api/internal/config"
	"gl-app/api/internal/svc"
	"gl-app/api/internal/worker"
)

var configFile = flag.String("f", "api/etc/gulu.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	serviceContext := svc.NewServiceContext(c)
	defer serviceContext.MediaQueue.Close()
	if err := worker.NewMediaWorker(ctx, serviceContext).Run(); err != nil && ctx.Err() == nil {
		logx.Must(err)
	}
}
