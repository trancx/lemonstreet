package http

import (
	"github.com/bilibili/kratos/pkg/ecode"
	"net/http"
	"search/internal/service"
	"strings"

	//v "search/api/vrfapi"

	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/log"
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
	"search/internal/model"
)

var (
	svc *service.Service
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
	svc = s
	engine = bm.DefaultServer(hc.Server)
	initRouter(engine)
	err = engine.Start()
	return
}

func initRouter(e *bm.Engine) {
	e.Ping(ping)
	e.GET("/search", search)
}

func search(c *bm.Context)  {
	var err error
	var result interface{}
	query := c.Request.URL.Query().Get("q")
	_type := c.Request.URL.Query().Get("type")

	// validate
	if strings.EqualFold(query, "") {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	// limiter

	// type filter
	switch _type {
	case "article":
		result, err = svc.SearchArticles(c, query)
	case "user":
		result, err = svc.SearchUsers(c, query)
	default:
		err = ecode.RequestErr
	}

	c.JSON(result, err)
}

func ping(ctx *bm.Context) {
	if _, err := svc.Ping(ctx, nil); err != nil {
		log.Error("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// example for http request handler.
func howToStart(c *bm.Context) {
	k := &model.Kratos{
		Hello: "Golang 大法好 !!!",
	}
	c.JSON(k, nil)
}