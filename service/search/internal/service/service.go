package service

import (
	"context"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/naming/discovery"
	"github.com/bilibili/kratos/pkg/net/rpc/warden/resolver"
	"search/api/accapi"
	"search/api/artapi"
	"search/internal/dao"

	"github.com/golang/protobuf/ptypes/empty"
)

// Service service.
type Service struct {
	ac  *paladin.Map
	dao dao.Dao
	accRPC api.AccountClient
	artRPC artapi.ArticleClient
}

// New new a service and return.
func New(d dao.Dao) (s *Service, err error) {
	var (
		accRPC api.AccountClient
		artRPC artapi.ArticleClient
	)
	if accRPC, err = api.NewRPCAccountClient(nil); err != nil {
		panic(err)
	}
	if artRPC, err = artapi.NewRPCArticleClient(nil); err != nil {
		panic(err)
	}
	s = &Service{
		ac:     &paladin.TOML{},
		dao:    d,
		accRPC: accRPC,
		artRPC: artRPC,
	}
	err = paladin.Watch("application.toml", s.ac)
	return
}

func init()  {
	resolver.Register(discovery.Builder())
}

func (s *Service) SearchArticles(c context.Context, query string) ([]*artapi.ArticleBaseInfo, error) {
	rpl, err := s.artRPC.SearchArticlesByTitle(c, &artapi.NameReq{
		Name:                 query,
	})
	if err != nil {
		return nil, err
	}
	return rpl.Infos, nil
}

func (s *Service) SearchUsers(c context.Context, query string) ([]*api.UserBaseInfo, error) {
	rpl, err := s.accRPC.SearchBaseInfoByName(c, &api.NameReq{
		Name:                 query,
	})
	if err != nil {
		return nil, err
	}
	return rpl.Infos, nil
}

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, s.dao.Ping(ctx)
}

// Close close the resource.
func (s *Service) Close() {
	s.dao.Close()
}
