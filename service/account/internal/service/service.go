package service

import (
	"context"
	"fmt"
	"kratos-demo/internal/model"

	pb "kratos-demo/api"
	"kratos-demo/internal/dao"
	"github.com/bilibili/kratos/pkg/conf/paladin"

	"github.com/golang/protobuf/ptypes/empty"
)

// Service service.
type Service struct {
	ac  *paladin.Map
	dao dao.Dao // interface, dao implement it !!
}

// New new a service and return
// Trance: 此处接受了 mysql mc redis 三者融合的一个 dao，并初始化了 service
func New(d dao.Dao) (s *Service, err error) {
	s = &Service{
		ac:  &paladin.TOML{},
		dao: d,
	}
	err = paladin.Watch("application.toml", s.ac)
	return
}

// SayHello grpc demo func.
func (s *Service) SayHello(ctx context.Context, req *pb.HelloReq) (reply *empty.Empty, err error) {
	reply = new(empty.Empty)
	fmt.Printf("hello %s", req.Name)
	return
}

// SayHelloURL bm demo func.
func (s *Service) SayHelloURL(ctx context.Context, req *pb.HelloReq) (reply *pb.HelloResp, err error) {
	reply = &pb.HelloResp{
		Content: "hello " + req.Name,
	}
	fmt.Printf("hello url %s", req.Name)
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

func (s *Service) Info(c context.Context, id int64) (res *model.UserInfo, err error) {
	_ = s.dao.UpdateAvatar(c, id, "old")
	res, err = s.dao.UserInfoID(c, id)
	//_ = s.dao.UpdateDesc(c, id, "update")
	//_ = s.dao.UpdateGender(c, id, "f")
	//_ = s.dao.UpdateMail(c, id, "update@qq.com")
	//res, err = s.dao.UserInfoID(c, id)
	return
}

func (s *Service) InfoName(c context.Context, name string) (res *model.UserInfo, err error) {
	res, err = s.dao.UserInfoName(c, name)
	return
}

func (s *Service) Account(c context.Context, info *model.UserInfo) (err error) {
	err = s.dao.SetUserInfo(c, info)
	return
}
