package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/log"
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
	pb "login/api"
	"login/internal/model"
)

const (
	_defaultDomain         = "localhost"
	_defaultCookieName     = "login_cookie"
	_defaultCookieLifeTime = 2592000

)

var svc pb.DemoServer

// New new a bm server.
func New(s pb.DemoServer) (engine *bm.Engine, err error) {
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
	pb.RegisterDemoBMServer(engine, s)
	initRouter(engine)
	err = engine.Start()
	return
}

func initRouter(e *bm.Engine) {
	e.Ping(ping)
	g := e.Group("/api")
	{
		g.POST("/login", login)
	}
}

func ping(ctx *bm.Context) {
	if _, err := svc.Ping(ctx, nil); err != nil {
		log.Error("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// example for http request handler.
func login(c *bm.Context) {
	info := &model.LoginInfo{}

	if err := c.Bind(&info); err != nil {
		fmt.Println("Bind error!")
	} else {
		fmt.Println("user " + info.Username + "login")
	}

	cookie := &http.Cookie{
		Name:     _defaultCookieName,
		Value:    "23333",
		Path:     "/",
		HttpOnly: true,
		Domain:   _defaultDomain,
	}
	cookie.MaxAge = _defaultCookieLifeTime
	cookie.Expires = time.Now().Add(time.Duration(_defaultCookieLifeTime) * time.Second)
	http.SetCookie(c.Writer, cookie)

	c.JSON(info, nil)
}
