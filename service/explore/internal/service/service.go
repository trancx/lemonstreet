package service

import (
	"context"
	"explore/internal/dao"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/naming/discovery"
	"github.com/bilibili/kratos/pkg/net/rpc/warden/resolver"

	"github.com/golang/protobuf/ptypes/empty"
	art "explore/api/artapi"
)

// Service service.
type Service struct {
	ac  *paladin.Map
	dao dao.Dao
	artRPC art.ArticleClient

}

func init()  {
	resolver.Register(discovery.Builder())
}

// New new a service and return.
func New(d dao.Dao) (s *Service, err error) {
	var (
		rpc art.ArticleClient
	)
	if rpc, err = art.NewRPCArticleClient(nil); err != nil {
		panic(err)
	}
	s = &Service{
		ac:  &paladin.TOML{},
		dao: d,
		artRPC:rpc,
	}
	err = paladin.Watch("application.toml", s.ac)
	return
}

func (s *Service) FetchArticles(c context.Context, beg int64, end int64) ([]*art.ArticleBaseInfo, error) {
	reply, err := s.artRPC.LatestArticles(c, &art.TimeReq{
		Beg:                  beg,
		End:                  end,
	})
	if err != nil {
		return nil, err
	}
	return reply.Infos, nil
}

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, s.dao.Ping(ctx)
}

// Close close the resource.
func (s *Service) Close() {
	s.dao.Close()
}
