package http

import (
	"fmt"
	"login/api/vrfapi"
	"login/internal/service"
	"net/http"
	"time"

	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/log"
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
	"login/internal/model"
)

const (
	_defaultDomain         = "localhost"
	_defaultCookieName     = "login_cookie"
	_defaultCookieLifeTime = 2592000

)

var (
	loginSvc *service.Service
	v		 *vrfapi.Verify
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
	loginSvc = s
	v = vrfapi.New()
	engine = bm.DefaultServer(hc.Server)
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

}

func genCookies(c *bm.Context, name string, value string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Domain:   _defaultDomain,
	}
	cookie.MaxAge = _defaultCookieLifeTime
	cookie.Expires = time.Now().Add(time.Duration(_defaultCookieLifeTime) * time.Second)
	http.SetCookie(c.Writer, cookie)
}

// example for http request handler.
func login(c *bm.Context) {
	var (
		uid int64
		err error
		token string
	)
	uid = 1
	info := model.LoginInfo{}

	if err := c.Bind(&info); err != nil {
		log.Error("Bind error! (%v)", err)
	} else {
		log.Info("user " + info.Username + "login")
	}

	if token, err = v.GenToken(c, uid); err != nil {
		c.JSON(nil, err)
		c.Abort()
		return
	}

	genCookies(c, "uid", fmt.Sprintf("%d", uid))
	genCookies(c, "token", token)

	c.JSON(info, nil)
}
