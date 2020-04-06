package http

import (
	"article/api/artapi"
	"article/internal/model"
	"article/internal/service"
	"fmt"
	"github.com/bilibili/kratos/pkg/ecode"
	"github.com/bilibili/kratos/pkg/net/http/blademaster/binding"
	"time"
	v "article/api/vrfapi"

	"github.com/bilibili/kratos/pkg/conf/paladin"
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
)

var (
	artSvc *service.Service
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
	artSvc = s
	verify = v.New()
	engine = bm.DefaultServer(hc.Server)
	initRouter(engine)
	err = engine.Start()
	return
}

func initRouter(e *bm.Engine) {
	//e.Ping(ping)
	g := e.Group("/api", verify.Verify)
	{
		g.POST("/article", postArticle)
		g.GET("/article", getArticle)
		g.DELETE("/article", delArticle)
		g.GET("/article/format", format)
	}

}

func test(c *bm.Context)  {
	c.Set("uid", int64(1))
}

func format(c *bm.Context) {
	var (
		apis []model.Format
		params struct{
			Aid int64	`json:"aid"`
		}
	)
	apis = append(apis, model.Format{
		Method: "get",
		API:    "/api/article",
		Params: &params,
	})
	apis = append(apis, model.Format{
		Method: "post",
		API:    "/api/article",
		Params: &model.PostArticle{},
	})

	c.JSON(apis, nil)
}

func delArticle(c *bm.Context)  {
	var (
		params struct{
			Aid int64	`json:"aid"`
		}
	)
	if err := c.BindWith(&params, binding.JSON); err != nil {
		return
	}
	uid, _ := c.Get("uid")
	c.JSON(nil, artSvc.DeleteArticle(c, params.Aid, uid.(int64)))
}

// 1. anounymous
// 2. authentic access
func getArticle(c *bm.Context) {
	var (
		params struct{
			Aid int64	`json:"aid"`
		}
	)
	if err := c.BindWith(&params, binding.JSON); err != nil {
		return
	}
	if res, err := artSvc.GetArticleAnnms(c, params.Aid); err != nil {
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
		desc []rune
		dlen int
	)
	uid, _ := c.Get("uid")
	pa = new(model.PostArticle)
	if err := c.BindWith(pa, binding.JSON); err != nil {
		return // aborted
	}
	desc = []rune(pa.Content)
	if dlen = len(desc); dlen > 50 {
		dlen = 50
	}
	// author fill by accRPC client
	abi = artapi.ArticleBaseInfo{
		Aid:    0,
		Uid:    uid.(int64),
		Title:  pa.Title,
		Desc:   fmt.Sprintf(string(desc[0:dlen])),
		Date:	time.Now().Unix(),
	}
	err := artSvc.PostArticle(c, &abi, pa.Content) // get user info
	if err != nil {
		c.JSON(nil, ecode.ServerErr)
		return
	}
	c.JSON(&abi, nil)
}