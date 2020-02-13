package http

import (
	"comment/internal/model"
	"comment/internal/service"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
)

var cmtSvc *service.Service

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
	cmtSvc = s
	engine = bm.DefaultServer(hc.Server)
	initRouter(engine)
	err = engine.Start()
	return
}

func initRouter(e *bm.Engine) {
	e.Ping(ping)
	//g := e.Group("/comment")
	//{
	//	g.GET("/start", howToStart)
	//}
	g := e.Group("/lemonstreet")
	{
		//g.GET("/:user/:title/comments") handle by RPC
		g.POST("/:user/:title/comments", postComment) // verify
	}
	e.GET("/format", format)
}

func ping(ctx *bm.Context) {

}

// example for http request handler.
func format(c *bm.Context) {
	c.JSON(model.PostComment{}, nil)
}

// example for http request handler.
func postComment(c *bm.Context) {
	cmm := model.PostComment{}
	if err := c.Bind(&cmm); err != nil {
		return
	}
	if err := cmtSvc.PostComment(c, &cmm.ABI, &cmm.Comment); err != nil {
		c.JSON(nil, err)
	} else {
		c.JSON(&cmm.Comment, nil)
	}
}