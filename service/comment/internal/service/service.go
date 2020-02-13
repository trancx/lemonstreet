package service

import (
	"comment/api/artapi"
	"comment/api/cmtapi"
	"comment/internal/dao"
	"context"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/ecode"
	"time"
)

// Service service.
type Service struct {
	ac  *paladin.Map
	dao dao.Dao
}

func (s *Service) CommentsOfAID(context.Context, *cmtapi.IDReq) (*cmtapi.CommentsReply, error) {
	panic("implement me")
}

func (s *Service) CommentSOfUID(context.Context, *cmtapi.IDReq) (*cmtapi.CommentsReply, error) {
	panic("implement me")
}

func (s *Service) PostComment(c context.Context, abi *artapi.ArticleBaseInfo,comment *cmtapi.Comment) error {
	var (
		err error
	)
	comment.Date = time.Now().Unix()
	err = s.dao.PostComment(c, comment)
	if err != nil {
		err = ecode.NothingFound
		return err
	}
	return nil
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
