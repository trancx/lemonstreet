package di

import (
	"context"
	"time"

	"verify/internal/service"

	"github.com/bilibili/kratos/pkg/log"
	"github.com/bilibili/kratos/pkg/net/rpc/warden"
)

//go:generate kratos tool wire
type App struct {
	svc *service.Service
	grpc *warden.Server
}

func NewApp(svc *service.Service, g *warden.Server) (app *App, closeFunc func(), err error){
	app = &App{
		svc: svc,
		grpc: g,
	}
	closeFunc = func() {
		ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
		if err := g.Shutdown(ctx); err != nil {
			log.Error("grpcSrv.Shutdown error(%v)", err)
		}
		svc.Close()
		cancel()
	}
	return
}