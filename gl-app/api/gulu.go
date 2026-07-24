package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpx"
	"gl-app/pkg/xcode"

	"gl-app/api/internal/config"
	"gl-app/api/internal/handler"
	"gl-app/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "api/etc/gulu.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	defer ctx.MediaQueue.Close()
	handler.RegisterHandlers(server, ctx)
	httpx.SetErrorHandler(xcode.ErrHandler)
	httpx.SetOkHandler(xcode.OkHandler) // 拦截器
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
