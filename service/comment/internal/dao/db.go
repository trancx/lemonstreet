package dao

import (
	"comment/api/cmtapi"
	"context"
	"github.com/prometheus/common/log"

	"comment/internal/model"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/database/sql"
)

const (
	_insertComment = "INSERT INTO `lemonstreet`.`comment` (`uid`, `aid`, `created`, `content`) VALUES (?,?,?,?)"
	_selCommentsByUId = "SELECT * FROM `comment` WHERE  `uid`=?"
	_selCommentsByAId = "SELECT * FROM `comment` WHERE  `aid`=?"
)

func NewDB() (db *sql.DB, err error) {
	var cfg struct {
		Client *sql.Config
	}
	if err = paladin.Get("db.toml").UnmarshalTOML(&cfg); err != nil {
		return
	}
	db = sql.NewMySQL(cfg.Client)
	return
}

func (d *dao) RawArticle(ctx context.Context, id int64) (art *model.Article, err error) {
	// get data from db
	return
}

// `uid`, `aid`, `created`, `content`
func (d *dao) PostComment(ctx context.Context, comment *cmtapi.Comment) (error) {
	result, err := d.db.Exec(ctx, _insertComment, comment.Uid, comment.Aid, comment.Date, comment.Content)
	if err != nil {
		log.Errorf("dao.PostComment Failed (%v)", err)
		return  err
	}
	cid, err := result.LastInsertId()
	if err != nil {
		log.Errorf("dao.PostComment GetLastInsertId Failed (%v)", err)
		return  err
	}
	comment.Cid = cid
	//d.cache.addCache  NEEDCACHE
	return nil
}

// NEEDCACHE
func (d *dao) SearchCommentsByUId(c context.Context, uid int64) (cmmts []*cmtapi.Comment, err error) {
	var (
		rows *sql.Rows
	)
	rows, err = d.db.Query(c, _selCommentsByUId, uid)
	for rows.Next() {
		temp := cmtapi.Comment{}
		err = rows.Scan(&temp.Cid, &temp.Uid, &temp.Aid, &temp.Date, &temp.Content)
		if err != nil {
			cmmts = nil
			return
		}
		cmmts = append(cmmts, &temp)
	}
	return
}

func (d *dao) SearchCommentsByAId(c context.Context, aid int64) (cmmts []*cmtapi.Comment, err error) {
	var (
		rows *sql.Rows
	)
	rows, err = d.db.Query(c, _selCommentsByAId, aid)
	for rows.Next() {
		temp := cmtapi.Comment{}
		err = rows.Scan(&temp.Cid, &temp.Uid, &temp.Aid, &temp.Date, &temp.Content)
		if err != nil {
			cmmts = nil
			return
		}
		cmmts = append(cmmts, &temp)
	}
	return
}