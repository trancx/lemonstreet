package grpc

import (
	"context"
	"fmt"
	"github.com/bilibili/kratos/pkg/conf/env"
	"github.com/bilibili/kratos/pkg/naming"
	"github.com/bilibili/kratos/pkg/naming/discovery"
	"github.com/trancx/netip/ip"
	"net"
	"os"
	"strings"
	pb "verify/api"

	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/net/rpc/warden"
)

const (
	discoveryID = "verify.service"
)

var (
	wardenIP net.IP
	wardenPort string
)

// New new a grpc server.
func New(svc pb.VerifyServer) (ws *warden.Server, cancel context.CancelFunc, err error) {
	var rc struct {
		Server *warden.ServerConfig
	}
	err = paladin.Get("grpc.toml").UnmarshalTOML(&rc)
	if err == paladin.ErrNotExist {
		err = nil
	}
	hn, _ := os.Hostname()
	dis := discovery.New(nil)
	if wardenIP, err = ip.GetExternalIP(); err != nil {
		panic(err)
	}
	wardenPort = strings.Split(rc.Server.Addr, ":")[1]
	fmt.Println("wadenpot = "+ wardenPort)
	ins := &naming.Instance {
		Zone:     env.Zone,
		Env:      env.DeployEnv,
		AppID:    discoveryID,
		Hostname: hn,
		Addrs: []string{
			"grpc://" + wardenIP.String() + ":" + wardenPort,
		},
	}
	cancel, err = dis.Register(context.Background(), ins)
	if err != nil {
		panic(err)
	}
	ws = warden.NewServer(rc.Server)
	pb.RegisterVerifyServer(ws.Server(), svc)
	ws, err = ws.Start()
	return
}
