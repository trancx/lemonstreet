package service

import (
	"account/internal/model"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/prometheus/common/log"

	pb "account/api"
	"account/internal/dao"
	"github.com/bilibili/kratos/pkg/conf/paladin"
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

func (s *Service) BaseInfoByName(c context.Context, req *pb.NameReq) (reply *pb.BaseInfoReply, err error) {
	res, err := s.dao.UserInfoName(c, req.Name)
	if err != nil {
		log.Errorf("rpc InfoByName faile (%v)", err)
		return
	}
	reply = &pb.BaseInfoReply{
		Info:	res.ToBaseInfo(),
	}
	return
}

func (s *Service) BaseInfosByName(c context.Context, req *pb.NamesReq) (reply *pb.BaseInfosReply, err error) {
	var (
		names 	[]string
		infos	[]*pb.UserBaseInfo
		temp	*model.UserInfo
	)
	names = req.Names
	infos = []*pb.UserBaseInfo{}
	for _, name := range names {
		temp, err = s.dao.UserInfoName(c, name)
		if err != nil {
			reply = nil
			log.Errorf("Service BaseInfosByName failed (%v)", err)
			return
		}
		infos = append(infos, temp.ToBaseInfo())
	}
	reply = &pb.BaseInfosReply{
		Infos:                infos,
	}
	return
}

func (s *Service) BaseInfo(c context.Context, req *pb.UidReq) (reply *pb.BaseInfoReply, err error) {
	res, err := s.dao.UserInfoID(c, req.Uid)
	if err != nil {
		log.Errorf("rpc InfoByName failed (%v)", err)
		return
	}
	reply = &pb.BaseInfoReply {
		Info:	res.ToBaseInfo(),
	}
	return
}

func (s *Service) SearchBaseInfoByName(c context.Context, req *pb.NameReq) (reply *pb.BaseInfosReply, err error) {
	var (
		name 	string
		infos  	[]*pb.UserBaseInfo
		raws	[]model.UserInfo
	)
	name = req.Name
	infos = []*pb.UserBaseInfo{}

	raws, err = s.dao.SearchUserInfoByName(c, name)
	if err != nil {
		log.Errorf("rpc SearchInfoByName failed (%v)", err)
		return
	}
	for _, temp := range raws {
		infos = append(infos, temp.ToBaseInfo())
	}
	reply = &pb.BaseInfosReply{
		Infos:                infos,
	}
	return
}

func (s *Service) CreateAccount(c context.Context, info *model.UserInfo) (err error) {
	err = s.dao.SetUserInfo(c, info)

	return
}

func (s *Service) UpdateAccount(c context.Context, info *model.UserInfo) (err error) {
	err = s.dao.SetUserInfo(c, info)

	// test case? avatar or name? or what?

	return
}

func (s *Service) UpdateAvatar(c context.Context, info *model.UserInfo) (err error) {
	err = s.dao.SetUserInfo(c, info)

	// test case? avatar or name? or what?

	return
}

func (s *Service) UpdateTel(c context.Context, info *model.UserInfo) (err error) {
	err = s.dao.SetUserInfo(c, info)

	// test case? avatar or name? or what?

	return
}

func (s *Service) UpdateName(c context.Context, info *model.UserInfo) (err error) {
	err = s.dao.SetUserInfo(c, info)

	// test case? avatar or name? or what?

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
