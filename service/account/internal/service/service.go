package service

import (
	pb "account/api/accapi"
	"account/internal/dao"
	"account/internal/model"
	"context"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/database/sql"
	"github.com/bilibili/kratos/pkg/naming/discovery"
	"github.com/bilibili/kratos/pkg/net/rpc/warden/resolver"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/prometheus/common/log"
)

// Service service.
type Service struct {
	ac  *paladin.Map
	dao *dao.Dao // interface, dao implement it !!
	// comment RPC
	//
}

func init() {
	resolver.Register(discovery.Builder())
}

// New new a service and return
// Trance: 此处接受了 mysql mc redis 三者融合的一个 dao，并初始化了 service
func New(d *dao.Dao) (s *Service, err error) {
	s = &Service{
		ac:  &paladin.TOML{},
		dao: d,
	}
	err = paladin.Watch("application.toml", s.ac)
	return
}

// tel -> DB -> Cached-ID
func (s *Service) BaseInfoByTel(c context.Context, req *pb.TelReq) (reply *pb.BaseInfoReply, err error) {
	var (
		tel string
		info *model.UserInfo
		uid int64
	)
	reply = &pb.BaseInfoReply{
		Initialize:true,
	}
	tel = req.Tel
	if info ,err = s.dao.UserInfoTel(c, tel); err != nil {
		if err != sql.ErrNoRows {
			reply = nil
			return
		}
		reply.Initialize = false; /* first time*/
		info = model.NewUser()
		info.Tel = tel
		if uid, err = s.dao.InsertUserInfo(c, info); err != nil {
			reply = nil
			return
		}
		info.UserID = uid
	}
	reply.Info = info.ToBaseInfo()
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

func (s *Service) Info(c context.Context, id int64) (info *model.UserInfo, err error) {
	info, err = s.dao.UserInfoID(c, id)
	if err != nil {
		log.Errorf("rpc InfoByName failed (%v)", err)
		return
	}
	return
}

func (s *Service) CreateAccount(c context.Context, info *model.UserInfo) (err error) {
	return
}

func (s *Service) UpdateAccount(c context.Context, info *model.UserInfo) (err error) {
	// test case? avatar or name? or what?
	return
}

func (s *Service) UpdateAvatar(c context.Context, uid int64, avatar string) (err error) {
	err = s.dao.UpdateAvatar(c, uid, avatar)
	// test case? avatar or name? or what?
	return
}

// sms rpc test
func (s *Service) UpdateMail(c context.Context, uid int64, mail string) (err error) {
	err = s.dao.UpdateMail(c, uid, mail)
	return
}

func (s *Service) UpdateGender(c context.Context, uid int64, gender string) (err error) {
	err = s.dao.UpdateGender(c, uid, gender)
	return
}

func (s *Service) UpdateDesc(c context.Context, uid int64, desc string) (err error) {
	err = s.dao.UpdateDesc(c, uid, desc)
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
