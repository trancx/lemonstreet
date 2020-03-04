package http

import (
	"github.com/bilibili/kratos/pkg/ecode"
	accapi "login/api/accapi"
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
	accRPC	accapi.AccountClient /* interface */
	v		 *vrfapi.Verify
)

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
	v = vrfapi.New() //panic itself
	if accRPC, err = accapi.NewRPCAccountClient(nil); err != nil {
		panic(err)
	}
	engine = bm.DefaultServer(hc.Server)
	initRouter(engine)
	err = engine.Start()
	return
}

func initRouter(e *bm.Engine) {
	g := e.Group("/api")
	{
		g.POST("/login", login)
		g.POST("/logout", logout)
		g.GET("/login/format", format)
	}
}

func format(c *bm.Context) {
	var (
		api []model.Format
	)
	api = append(api, model.Format{
		Method: "post",
		API: "/api/login",
		Params: &model.LoginInfo{},
	})
	api = append(api, model.Format{
		Method: "post",
		API: "/api/loout",
		Params: nil,
	})
	c.JSON(api, nil)
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

func login(c *bm.Context) {
	var (
		uid int64
		err error
		token string
	)
	params := model.LoginInfo{}

	if err := c.Bind(&params); err != nil {
		log.Error("Bind error! (%v)", err)
	} else {
		log.Info("user " + params.Tel + "login")
	}

	// sms check smsRPC.verify(...)

	// user rpc
	reply, err := accRPC.BaseInfoByTel(c, &accapi.TelReq{
		Tel:                  params.Tel,
		RealIp:               "",
	})

	if err != nil {
		c.JSON(nil, ecode.ServerErr)
		return
	}
	// first time or not? ecode to reprent the status
	uid = reply.Info.Uid
	if token, err = v.GenToken(c, uid); err != nil {
		c.JSON(nil, err)
		return
	}
	genCookies(c, _defaultCookieName, fmt.Sprintf("uid=%d&token=%s", uid, token))

	c.JSON(reply.Info, nil)
}

func logout(c *bm.Context) {
 	// del cookie
}