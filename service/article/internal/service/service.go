package service

import (
	acc "article/api/accapi"
	artapi "article/api/artapi"
	"article/internal/dao"
	"article/internal/model"
	"context"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/ecode"
)

// Service service.
type Service struct {
	ac  *paladin.Map
	dao dao.Dao
}

/* gRPC implementation, normal path, db -> cache, or specific amount(e.g. 100) comments
	cache together? but we have both uid-key and aid-key, cache is foolish here.
*/

func (s *Service) SearchArticlesByUID(context.Context, *artapi.IDReq) (*artapi.ArticleBaseInfosReply, error) {
	panic("implement me")
}

func (s *Service) SearchArticlesByTitle(context.Context, *artapi.NameReq) (*artapi.ArticleBaseInfosReply, error) {
	panic("implement me")
}

// username && title
// content and userinfo inside http body!
func (s *Service) PostArticle(c context.Context, info *artapi.ArticleBaseInfo, content string) (err error) {
	// save to DB and cache it!
	err = s.dao.PostArticle(c, info, content)
	// handle error
	return
}

// generate comments
func (s *Service) GetArticleAnnms(c context.Context, uname string, title string) (*model.ArticleInfo, error) {
	var (
		reply *acc.BaseInfoReply
		cdds  []artapi.ArticleBaseInfo
		abi	*artapi.ArticleBaseInfo
		art *model.Article
		info *acc.UserBaseInfo
		err	error
	)
	reply, err = s.dao.UserBaseInfoByName(c, uname)
	if err != nil {
		err = ecode.NothingFound
		return nil, err
	}
	info = reply.Info
	cdds, err = s.dao.ArticleBaseInfosByName(c, title)  // cache it!
	if err != nil {
		err = ecode.NothingFound
		return nil, err
	}

	for _, temp := range cdds {
		if temp.Uid == info.Uid {
			abi = &temp	// weird
			break
		}
	}
	if abi == nil {
		err = ecode.NothingFound
		return nil, err
	}
	art, err = s.dao.Article(c, abi.Aid)
	if err != nil {
		err = ecode.ServerErr
		return nil, err
	}
	return &model.ArticleInfo{
		UInfo:   info,
		ABI:     abi,
		Content: art,
	}, err
}

func (s *Service) GetArticle(c context.Context, user *acc.UserBaseInfo, article *model.ArticleBaseInfo) (info *model.ArticleInfo, err error) {
	var (

	)

	// generate article_id and fill the article baseinfo

	// save to DB and cache it!

	return
}

// New new a service and return.
func New(d dao.Dao) (s *Service, err error) {
	s = &Service{
		ac:  &paladin.TOML{},
		dao: d,
	}
	err = paladin.Watch("application.toml", s.ac)
	return
}

// Close close the resource.
func (s *Service) Close() {
	s.dao.Close()
}
