package http

import (
	"account/internal/model"
	"fmt"
	"net/http"

	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/log"
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"

	"account/internal/service"
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

// domain/[:username]
func initRouter(e *bm.Engine) {
	e.Ping(ping)
	//g := e.Group("/api/account")
	//{
	//	g.GET("/info", info)
	//	g.GET("infobyname", infoName)
	//	g.GET("/search", search)
	//}
	e.GET("/api/:user", test)
}

// example for http request handler.
func test(c *bm.Context) {
	res := c.Request.URL.Query().Get("user")
	reply := model.UserInfo{
		Name:	res,
	}

	c.JSON(&reply, nil)
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
	res, err := accSvc.Info1(c,27182818285)

	if err != nil {
		fmt.Println("%v!", err)
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

func search(c *bm.Context) {
	var (
		q string
	)

	q = c.Request.URL.Query().Get("q")
	// 解析 json -> go-model -> dao -> context
	res, err := accSvc.Search(c,q)

	if err != nil {
		fmt.Println("error!")
	}


	fmt.Println(len(res))


	c.JSON(res, nil)
}