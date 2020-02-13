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
