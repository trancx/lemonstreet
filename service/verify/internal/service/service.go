package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/prometheus/common/log"
	"strings"
	"time"
	vrfapi "verify/api"
	"verify/internal/dao"
)

// Service service.
type Service struct {
	ac  *paladin.Map
	dao dao.Dao
}

func (s *Service) UdtKey(c context.Context, req *vrfapi.TokenReq) (*vrfapi.TokenReply, error) {
	now := time.Now().Nanosecond() + int(req.Tk.Id)
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%d", now)))
	token := hash.Sum(nil)
	req.Tk.Key = fmt.Sprintf("%x", token)

	if err := s.dao.UpdateKey(c, req.Tk); err != nil {
		return nil, err
	}
	log.Info("GenKey uid: %d token: %s", req.Tk.Id, req.Tk.Key)
	return &vrfapi.TokenReply{
		IsValid:              true,
		Tk:                   req.Tk,
	}, nil
}

func (s *Service) GenKey(c context.Context, req *vrfapi.TokenReq) (*vrfapi.TokenReply, error) {
	now := time.Now().Nanosecond() + int(req.Tk.Id)
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%d", now)))
	token := hash.Sum(nil)
	req.Tk.Key = fmt.Sprintf("%x", token)

	// 当用户清除 cookies 就会到这个里面, 当然判断用户是否是第一次登陆也可以保证
	if err := s.dao.InsertKey(c, req.Tk); err != nil {
		log.Infof("GenKey Failed, Try UpdateKey uid: %d token: %s", req.Tk.Id, req.Tk.Key)
		if err = s.dao.UpdateKey(c, req.Tk); err != nil {
			log.Errorf("UpdateKey Failed (%d)", err)
			return nil, err
		}
	}
	log.Infof("GenKey uid: %d token: %s", req.Tk.Id, req.Tk.Key)
	return &vrfapi.TokenReply{
		IsValid:              true,
		Tk:                   req.Tk,
	}, nil
}

func (s *Service) VrfKey(c context.Context, req *vrfapi.TokenReq) (reply *vrfapi.TokenReply, err error) {
	var (
		token *vrfapi.Token
	)
	reply = &vrfapi.TokenReply{
		IsUpdated: 			  false,
		IsValid:              false,
		Tk:                   nil,
	}
	token, err = s.dao.GetKey(c, req.Tk.Id)
	if err != nil {
		return
	}
	if !strings.EqualFold(token.Key, req.Tk.Key) {
		reply.IsUpdated = true
		reply.Tk = token
		return
	}
	reply.IsValid = true
	return
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

// Close close the resource.
func (s *Service) Close() {
	s.dao.Close()
}
