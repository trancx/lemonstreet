package grpc

import (
	pb "article/api/artapi"

	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/net/rpc/warden"
)

// New new a grpc server.
func New(svc pb.ArticleServer) (ws *warden.Server, err error) {
	var rc struct {
		Server *warden.ServerConfig
	}
	err = paladin.Get("grpc.toml").UnmarshalTOML(&rc)
	if err == paladin.ErrNotExist {
		err = nil
	}
	ws = warden.NewServer(rc.Server)
	pb.RegisterArticleServer(ws.Server(), svc)
	ws, err = ws.Start()
	return
}
