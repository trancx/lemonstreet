package service

import (
	"context"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/naming/discovery"
	"github.com/bilibili/kratos/pkg/net/rpc/warden/resolver"
	"login/api"
	"login/internal/dao"

	"github.com/golang/protobuf/ptypes/empty"
)

// Service service.
type Service struct {
	ac  *paladin.Map
	dao dao.Dao
}

func (s *Service) SayHello(context.Context, *api.HelloReq) (*empty.Empty, error) {
	panic("implement me")
}

func (s *Service) SayHelloURL(context.Context, *api.HelloReq) (*api.HelloResp, error) {
	panic("implement me")
}

func init() {
	// NOTE: 注意这段代码，表示要使用discovery进行服务发现
	// NOTE: 还需注意的是，resolver.Register是全局生效的，所以建议该代码放在进程初始化的时候执行
	// NOTE: ！！！切记不要在一个进程内进行多个不同中间件的Register！！！
	// NOTE: 在启动应用时，可以通过flag(-discovery.nodes) 或者 环境配置(DISCOVERY_NODES)指定discovery节点
	resolver.Register(discovery.Builder())
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

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, s.dao.Ping(ctx)
}

// Close close the resource.
func (s *Service) Close() {
	s.dao.Close()
}
