package dao

import (
	"context"
	"github.com/prometheus/common/log"
	vrfapi "verify/api"

	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/database/sql"
)

const (
	_insertKey = "INSERT INTO `lemonstreet`.`token` (`id`, `token`) VALUES (?,?)"
	_updateKey = "UPDATE `lemonstreet`.`token` SET `token` = ? WHERE (`id` = ?)"
	_selKey = "SELECT * FROM `lemonstreet`.`token` WHERE (`id` = ?)"
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

func (d *dao) RawGetKey(ctx context.Context, id int64) (token *vrfapi.Token, err error) {
	token = &vrfapi.Token{}
	err = d.db.QueryRow(ctx, _selKey, id).Scan(&token.Id, &token.Key)
	if err != nil {
		log.Errorf("dao.RawGetKey Failed (%d)", err)
	}
	return
}

func (d *dao) InsertKey(c context.Context, key *vrfapi.Token) error {
	var (
		err error
	)
	_, err = d.db.Exec(c, _insertKey, key.Id, key.Key)
	if err != nil {
		log.Errorf("dao.InsertKey Failed (%d)", err)
	}
	return err
}

func (d *dao) UpdateKey(c context.Context, key *vrfapi.Token) error {
	var (
		err error
	)
	_, err = d.db.Exec(c, _updateKey, key.Key, key.Id)
	if err != nil {
		log.Errorf("dao.UpdateKey Failed (%d)", err)
	}
	return err
}
