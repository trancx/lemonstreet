package http

import (
	"github.com/bilibili/kratos/pkg/ecode"
	"net/http"
	"search/internal/model"
	"search/internal/service"
	"strings"

	//v "search/api/vrfapi"

	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/log"
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
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
	g := e.Group("/api")
	{
		g.GET("/search", search)
		g.GET("/search/format", format)
	}

}

func format(c *bm.Context)  {
	var (
		api []model.Format
	)
	api = append(api, model.Format{
		Method: "get",
		API: "/api/search?q=?&type=?",
		Params: "q=?(string)&type=?(article/user)",
	})
	c.JSON(api, nil)
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