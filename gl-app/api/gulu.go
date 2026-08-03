package main

import (
	"context"
	"flag"
	"fmt"

	"gl-app/api/internal/config"
	"gl-app/api/internal/handler"
	"gl-app/api/internal/svc"
	"gl-app/api/internal/worker"
	"gl-app/pkg/xcode"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "api/etc/gulu.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)

	ctx := svc.NewServiceContext(c)
	defer ctx.Close()
	handler.RegisterHandlers(server, ctx)
	httpx.SetErrorHandler(xcode.ErrHandler)
	httpx.SetOkHandler(xcode.OkHandler) // 拦截器

	// HTTP 服务与Kafka媒体消费者使用同一生命周期，启动系统即可处理转码任务。
	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()
	serviceGroup.Add(server)
	serviceGroup.Add(worker.NewMediaWorker(context.Background(), ctx))

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	serviceGroup.Start()
}
