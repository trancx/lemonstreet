package http

import (
	"github.com/bilibili/kratos/pkg/ecode"
	accapi "login/api/accapi"
	"fmt"
	"login/api/vrfapi"
	"login/internal/service"
	"net/http"
	"strings"
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
	cms = map[int]string{
		0: "ok",
		-601: "Sms check fail",
		-602: "Too many times",
		-603: "User dosn't login",
	}
)

func initEcode()  {
	ecode.Register(cms)
}

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
		g.POST("/logout", v.Verify, logout)  // verify middleware FIXME
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
		Errs:cms,
	})
	api = append(api, model.Format{
		Method: "post",
		API: "/api/loout",
		Params: nil,
		Errs:cms,
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

func delCookies(c *bm.Context, name string) {
	var (
		cookie *http.Cookie
		err error
	)
	if cookie, err = c.Request.Cookie(name); err != nil {
		return
	}
	cookie.Value = "nil"
	cookie.MaxAge = 0
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
	if len(params.Tel) != 11 || len(params.Sms) != 6 || !strings.HasSuffix(params.Tel, params.Sms) {
		c.JSON(nil, ecode.Int(-601))
		return
	}

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
	var (
		err error
	)
	uid, _ := c.Get("uid")
	if _, err = v.GenToken(c, uid.(int64)); err != nil {
		log.Error("Token invalid (%v)", err)
		c.JSON(nil, err)
		return
	}
	delCookies(c, _defaultCookieName)
	c.JSON(nil, nil)
}