// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package di

import (
	pb "verify/api"
	"verify/internal/dao"
	"verify/internal/server/grpc"
	"verify/internal/server/http"
	"verify/internal/service"

	"github.com/google/wire"
)

var daoProvider = wire.NewSet(dao.New, dao.NewDB, dao.NewRedis, dao.NewMC)
var serviceProvider = wire.NewSet(service.New, wire.Bind(new(pb.DemoServer), new(*service.Service)))

func InitApp() (*App, func(), error) {
	panic(wire.Build(daoProvider, serviceProvider, http.New, grpc.New, NewApp))
}
