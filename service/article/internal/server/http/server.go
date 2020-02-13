package http

import (
	"article/api/artapi"
	"article/internal/model"
	"article/internal/service"
	"fmt"
	"github.com/bilibili/kratos/pkg/ecode"
	"time"

	"github.com/bilibili/kratos/pkg/conf/paladin"
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
)

var artSvc *service.Service

func initRouter(e *bm.Engine) {
	//e.Ping(ping)
	g := e.Group("/lemonstreet")
	{
		g.POST("/:user/:title", postArticle)
		g.GET("/:user/:title", getArticle)

	}
	e.GET("/format", format)
}

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
	artSvc = s
	engine = bm.DefaultServer(hc.Server)
	initRouter(engine)
	err = engine.Start()
	return
}

func format(c *bm.Context) {
	content := new(model.PostArticle)
	c.JSON(content, nil)
}

// 1. anounymous
// 2. authentic access
func getArticle(c *bm.Context) {
	user, _ := c.Params.Get("user")
	title, _:= c.Params.Get("title")

	if res, err := artSvc.GetArticleAnnms(c, user, title); err != nil {
		c.JSON(res, err)
	} else {
		c.JSON(res, nil)
	}
}

// authentic request.  body contains userbaseinfo & articlebaseinfo
// make sure uid && uname & logger valid !!
func postArticle(c *bm.Context) {
	var (
		pa    *model.PostArticle
		abi   artapi.ArticleBaseInfo
		title string
		desc []rune
		dlen int
	)

	pa = new(model.PostArticle)
	if err := c.Bind(pa); err != nil {
		return // aborted
	}
	desc = []rune(pa.Content)
	if dlen = len(desc); dlen > 50 {
		dlen = 50
	}

	title, _ = c.Params.Get("title")
	fmt.Println("title: " + title )
	abi = artapi.ArticleBaseInfo{
		Aid:    0,
		Uid:    pa.UBaseInfo.Uid,
		Author: pa.UBaseInfo.Name,
		Title:  title,
		Desc:   fmt.Sprintf(string(desc[0:dlen])),
		Date:	time.Now().Unix(),
	}
	err := artSvc.PostArticle(c, &abi, pa.Content) // get user info

	if err != nil {
		c.JSON(nil, ecode.ServerErr)
		return
	}
	c.JSON(&abi, err)
}