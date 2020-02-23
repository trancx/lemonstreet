package http

import (
	"account/internal/model"
	"fmt"
	"github.com/bilibili/kratos/pkg/ecode"
	"github.com/bilibili/kratos/pkg/net/http/blademaster/binding"
	"net/http"

	v "account/api/vrfapi"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/log"
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"

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
	//verify = v.New()
	engine = bm.DefaultServer(hc.Server)
	// pb.RegisterDemoBMServer(engine, s) 这里可以算是适配 protoc 的中间层
	initRouter(engine)
	err = engine.Start()
	return
}

func test(c *bm.Context) {
	c.Set("uid", int64(1))
}

// domain/[:username]
func initRouter(e *bm.Engine) {
	e.Ping(ping)
	g := e.Group("/api/account", test)	// FIXME  tourist or not? , verify.Verify
	{
		g.GET("/info", getUserInfo)

		g.POST("/avatar", postAvatar)
		g.POST("/email", postMail)
		g.POST("/desc", postDesc)
		g.POST("/gender", postGender)
	}
}

// example for http request handler.
func getUserInfo(c *bm.Context) {
	var (
		info struct{
			UId int64  	`json:"uid"`
			Name string `json:"name"`
		}
	)
	uid, exist := c.Get("uid")

	if err := c.BindWith(&info, binding.JSON); err != nil {
		return
	}
	// if uid == uid, can moreinfo about tel & other
	reply, err := accSvc.Info(c, info.UId)
	// FIXME: 404 not found. status transfer- sql -> http or gRPC
	if err != nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}

	if !exist || uid != info.UId {
		reply.Tel = ""
		reply.Mail = ""
	}

	c.JSON(reply, nil)
}

func postMail(c *bm.Context) {
	var (
		mail struct {
			Mail string `json:"mail"`
		}
	)
	if err := c.BindWith(&mail, binding.JSON); err != nil {
		return
	}
	uid, _ := c.Get("uid")

	if err := accSvc.UpdateMail(c, uid.(int64), mail.Mail); err != nil {
		err = ecode.ServerErr
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func postDesc(c *bm.Context) {
	var (
		desc struct {
			Desc string `json:"desc"`
		}
	)
	if err := c.BindWith(&desc, binding.JSON); err != nil {
		return
	}
	uid, _ := c.Get("uid")

	if err := accSvc.UpdateDesc(c, uid.(int64), desc.Desc); err != nil {
		err = ecode.ServerErr
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func postGender(c *bm.Context) {
	var (
		gender struct {
			Gender string `json:"gender"`
		}
	)
	if err := c.BindWith(&gender, binding.JSON); err != nil {
		return
	}
	uid, _ := c.Get("uid")

	if err := accSvc.UpdateGender(c, uid.(int64), gender.Gender); err != nil {
		err = ecode.ServerErr
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func postAvatar(c *bm.Context) {
	var (
		avatar struct {
			Avatar string `json:"avatar"`
		}
	)
	if err := c.BindWith(&avatar, binding.JSON); err != nil {
		return
	}
	uid, _ := c.Get("uid")

	if err := accSvc.UpdateAvatar(c, uid.(int64), avatar.Avatar); err != nil {
		err = ecode.ServerErr
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
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