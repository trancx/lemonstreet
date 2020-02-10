package service

import (
	artapi "article/api/artapi"
	acc		"article/api/accapi"
	"article/internal/dao"
	"context"
	"github.com/bilibili/kratos/pkg/conf/paladin"

	"github.com/golang/protobuf/ptypes/empty"
)

// Service service.
type Service struct {
	ac  *paladin.Map
	dao dao.Dao
}

func (s *Service) SearchArticlesByUID(context.Context, *artapi.IDReq) (*artapi.ArticleBaseInfosReply, error) {
	panic("implement me")
}

func (s *Service) SearchArticlesByTitle(context.Context, *artapi.NameReq) (*artapi.ArticleBaseInfosReply, error) {
	panic("implement me")
}

func (s *Service) Content(c context.Context, name string) (reply *acc.BaseInfoReply, err error) {

	return s.dao.Content(c, name)
}

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, s.dao.Ping(ctx)
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
