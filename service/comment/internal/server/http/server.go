package http

import (
	cmtapi "comment/api/cmtapi"
	"comment/internal/model"
	"comment/internal/service"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/ecode"
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
	v "comment/api/vrfapi"
	"github.com/bilibili/kratos/pkg/net/http/blademaster/binding"
)

var (
	cmtSvc *service.Service
	verify *v.Verify
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
	cmtSvc = s
	verify = v.New()
	engine = bm.DefaultServer(hc.Server)
	initRouter(engine)
	err = engine.Start()
	return
}

func initEcode() {
	// FIXME: ecode register
}

func initRouter(e *bm.Engine) {
	e.Ping(ping)
	//g := e.Group("/comment")
	//{
	//	g.GET("/start", howToStart)
	//}
	g := e.Group("/api", verify.Verify)
	{
		//g.GET("/:user/:title/comments") handle by RPC
		g.POST("/comment", postComment) // verify
		g.GET("/comment/format", format)
	}

}

func test(c *bm.Context)  {
	c.Set("uid", int64(1))
}

func ping(ctx *bm.Context) {

}

// example for http request handler.
func format(c *bm.Context) {
	var (
		api []model.Format
		params struct{
			AId int64 `json:"aid"`
			Content string `json:"content"`
		}
	)
	api = append(api, model.Format{
		Method: "post",
		API: "/api/comment",
		Params: &params,
	})
	c.JSON(api, nil)
}

// example for http request handler.
func postComment(c *bm.Context) {
	var (
		params struct{
			AId int64 `json:"aid"`
			Content string `json:"content"`
		}
		cmm = &cmtapi.Comment{}
	)

	if err := c.BindWith(&params, binding.JSON); err != nil {
		return
	}

	if len(params.Content) < 15 {
		c.JSON(nil, ecode.RequestErr) // error too short!! FIXME
		return
	}

	uid, _ := c.Get("uid")
	cmm.Content = params.Content
	cmm.Aid = params.AId
	cmm.Uid = uid.(int64)
	// FIXME: aid validation!!!
	// and it can delay!
	if err := cmtSvc.PostComment(c, cmm); err != nil {
		c.JSON(nil, err)
	} else {
		c.JSON(nil, nil)
	}
}