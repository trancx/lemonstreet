package dao

import (
	"context"

	"search/internal/model"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/database/sql"
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
