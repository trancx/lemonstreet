package http

import (
	api "account/api/accapi"
	"account/internal/model"
	"fmt"
	"github.com/bilibili/kratos/pkg/ecode"
	"net/http"

	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/log"
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
	v "account/api/vrfapi"

	"account/internal/service"
)

//var accSvc pb.DemoServer

var (
	accSvc *service.Service
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
	accSvc = s
	verify = v.New()
	engine = bm.DefaultServer(hc.Server)
	// pb.RegisterDemoBMServer(engine, s) 这里可以算是适配 protoc 的中间层
	initRouter(engine)
	err = engine.Start()
	return
}

// domain/[:username]
func initRouter(e *bm.Engine) {
	e.Ping(ping)
	g := e.Group("/lemonstreet", verify.Verify)	// FIXME  tourist or not?
	{
		g.GET("/:user", getUserInfo)

		g.POST("/:user/avatar", postUserInfo)
		g.POST("/:user/name", postUserInfo)
		g.POST("/:user/telephone", postUserInfo)
		g.POST("/:user", postUserInfo)
	}
}

// example for http request handler.
func getUserInfo(c *bm.Context) {

	userName, _ := c.Params.Get("user")
	//uid, _ := c.Get("uid")

	reply, err := accSvc.BaseInfoByName(c, &api.NameReq{
		Name:                 userName,
	})
	// FIXME: 404 not found. status transfer- sql -> http or gRPC
	if err != nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	c.JSON(reply.Info, nil)
}

// new user
func postUserInfo(c *bm.Context) {
	var (
		reply 	*model.UserInfo
		err		error
	)
	userName, _ := c.Params.Get("user")

	fmt.Println(userName)

	if err != nil {

	}

	c.JSON(reply, nil)
}

func updateUserInfo(c *bm.Context) {
	var (
		reply 	*model.UserInfo
		err		error
	)
	userName, _ := c.Params.Get("user")

	fmt.Println(userName)
	if err != nil {

	}
	c.JSON(reply, nil)
}


func ping(ctx *bm.Context) {
	if _, err := accSvc.Ping(ctx, nil); err != nil {
		log.Error("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}