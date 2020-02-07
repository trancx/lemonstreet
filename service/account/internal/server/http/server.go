package http

import (
	"fmt"
	"net/http"

	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/log"
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"

	"kratos-demo/internal/service"
)

//var accSvc pb.DemoServer

var (
	accSvc *service.Service
	//verify *v.Verify
)

// New new a bm server.
func New(s *service.Service) (engine *bm.Engine, err error) {
	var (
		hc struct {
			Server *bm.ServerConfig
		}
	)
	if err = paladin.Get("http.toml").UnmarshalTOML(&hc); err != nil {
		if err != paladin.ErrNotExist {
			return
		}
		err = nil
	}
	accSvc = s
	engine = bm.DefaultServer(hc.Server)
	// pb.RegisterDemoBMServer(engine, s) 这里可以算是适配 protoc 的中间层
	initRouter(engine)
	err = engine.Start()
	return
}

func initRouter(e *bm.Engine) {
	e.Ping(ping)
	g := e.Group("/api/account")
	{
		g.GET("/info", info)
		g.GET("infobyname", infoName)
	}
}

func ping(ctx *bm.Context) {
	if _, err := accSvc.Ping(ctx, nil); err != nil {
		log.Error("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// example for http request handler.
func info(c *bm.Context) {
	// 解析 json -> go-model -> dao -> context
	res, err := accSvc.Info(c,27182818285)

	if err != nil {
		fmt.Println("error!")
	}

	c.JSON(res, nil)
}

func infoName(c *bm.Context) {
	// 解析 json -> go-model -> dao -> context
	res, err := accSvc.InfoName(c,"trance")

	if err != nil {
		fmt.Println("error!")
	}

	c.JSON(res, nil)
}