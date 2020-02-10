package http

import (
	"article/internal/service"
	"fmt"
	"net/http"

	"article/internal/model"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/log"
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
)

var artSvc *service.Service

func initRouter(e *bm.Engine) {
	e.Ping(ping)
	g := e.Group("/lemonstreet")
	{
		g.GET("/:user/:title", article)
		g.POST("/:user/:title", howToStart)
	}
}

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
	artSvc = s
	engine = bm.DefaultServer(hc.Server)
	initRouter(engine)
	err = engine.Start()
	return
}

func ping(ctx *bm.Context) {
	if _, err := artSvc.Ping(ctx, nil); err != nil {
		log.Error("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// example for http request handler.
func article(c *bm.Context) {
	uname, _ := c.Params.Get("user")
	title, _ := c.Params.Get("title")

	fmt.Println("title: " + title)
	reply ,_:= artSvc.Content(c, uname)

	c.JSON(reply.Info, nil)
}

// example for http request handler.
func howToStart(c *bm.Context) {
	k := &model.Kratos{
		Hello: "Golang 大法好 !!!",
	}
	c.JSON(k, nil)
}