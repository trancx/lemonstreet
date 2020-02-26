package service

import (
	"comment/api/cmtapi"
	"comment/internal/dao"
	"context"
	"database/sql"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/ecode"
	"time"
)

// Service service.
type Service struct {
	ac  *paladin.Map
	dao dao.Dao
}

func (s *Service) CommentsOfAID(c context.Context, req *cmtapi.IDReq) (*cmtapi.CommentsReply, error) {
	res, err := s.dao.SearchCommentsByAId(c, req.Id)
	if err != nil && err != sql.ErrNoRows {
		err = ecode.NothingFound
		return nil, err
	}

	return &cmtapi.CommentsReply{
		Comments:             res,
	}, nil
}

func (s *Service) CommentSOfUID(c context.Context, req *cmtapi.IDReq) (*cmtapi.CommentsReply, error) {
	res, err := s.dao.SearchCommentsByUId(c, req.Id)
	if err != nil && err != sql.ErrNoRows {
		err = ecode.NothingFound
		return nil, err
	}

	return &cmtapi.CommentsReply{
		Comments:             res,
	}, nil
}

func (s *Service) PostComment(c context.Context, comment *cmtapi.Comment) error {
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
