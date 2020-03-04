// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package di

import (
	"github.com/google/wire"
	"article/internal/dao"
	"article/internal/server/grpc"
	"article/internal/server/http"
	"article/internal/service"
)

// Injectors from wire.go:

func InitApp() (*App, func(), error) {
	redis, err := dao.NewRedis()
	if err != nil {
		return nil, nil, err
	}
	memcache, err := dao.NewMC()
	if err != nil {
		return nil, nil, err
	}
	db, err := dao.NewDB()
	if err != nil {
		return nil, nil, err
	}
	daoDao, err := dao.New(redis, memcache, db)
	if err != nil {
		return nil, nil, err
	}
	serviceService, err := service.New(daoDao)
	if err != nil {
		return nil, nil, err
	}
	engine, err := http.New(serviceService)
	if err != nil {
		return nil, nil, err
	}
	server, cancel, err := grpc.New(serviceService)
	if err != nil {
		return nil, nil, err
	}
	app, cleanup, err := NewApp(serviceService, engine, server)
	if err != nil {
		return nil, nil, err
	}
	return app, func() {
		cancel()
		cleanup()
	}, nil
}

// wire.go:

var daoProvider = wire.NewSet(dao.New, dao.NewDB, dao.NewRedis, dao.NewMC)

