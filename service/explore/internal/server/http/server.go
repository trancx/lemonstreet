package http

import (
	"explore/api/artapi"
	"explore/internal/model"
	"explore/internal/service"
	"github.com/bilibili/kratos/pkg/ecode"
	"github.com/bilibili/kratos/pkg/net/http/blademaster/binding"
	"net/http"

	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/log"
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
)

var (
	svc *service.Service
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
	svc = s
	engine = bm.DefaultServer(hc.Server)
	initRouter(engine)
	err = engine.Start()
	return
}

func initRouter(e *bm.Engine) {
	e.Ping(ping)
	e.GET("/explore", explore)
	e.GET("/api/explore/format", format)
}

func ping(ctx *bm.Context) {
	if _, err := svc.Ping(ctx, nil); err != nil {
		log.Error("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func format(c *bm.Context)  {
	var (
		period struct {
			Begin int64 `json:"begin"`
			End int64 	`json:"end"`
		}
	)
	apis := model.Format{
		Method: "get",
		API:    "/explore",
		Params: &period,
	}
	c.JSON(apis, nil)
}

// example for http request handler.
func explore(c *bm.Context) {
	var (
		period struct {
			Begin int64 `json:"begin"`
			End int64 	`json:"end"`
		}
		// timestamp before begin or after end
		abis []*artapi.ArticleBaseInfo
	)
	if err := c.BindWith(&period, binding.JSON); err != nil {
		return
	}
	abis, err := svc.FetchArticles(c, period.Begin, period.End)
	if err != nil {
		c.JSON(nil, ecode.ServerErr)
		return
	}
	c.JSON(abis, nil)
}