// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package di

import (
	"verify/internal/dao"
	"verify/internal/server/grpc"
	"verify/internal/service"
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
	server, err := grpc.New(serviceService)
	if err != nil {
		return nil, nil, err
	}
	app, cleanup, err := NewApp(serviceService, server)
	if err != nil {
		return nil, nil, err
	}
	return app, func() {
		cleanup()
	}, nil
}

// wire.go:


