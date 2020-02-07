package dao

import (
	"context"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/database/sql"
	"github.com/bilibili/kratos/pkg/log"
	"kratos-demo/internal/model"
)

const (
	_selUserInfoID = "SELECT uid,name,tel,mail,gender,avatar,description,created FROM user WHERE uid=? "
	_selUserInfoName = "SELECT uid,name,tel,mail,gender,avatar,description,created FROM user WHERE name=?"

	)

func NewDB() (db *sql.DB, err error) {
	var cfg struct {
		Client *sql.Config
	}
	if err = paladin.Get("db.toml").UnmarshalTOML(&cfg); err != nil {
		return
	}
	//fmt.Println(cfg.Client)
	//info := model.UserInfo{}
	//id := 27182818285
	db = sql.NewMySQL(cfg.Client)
	err = db.Ping(context.Background())
	if err != nil {
		panic(err)
	}
	//err = db.QueryRow(nil, _selUserInfoID, id).Scan(&info.UserID, &info.Name, &info.Mail, &info.Gender, &info.Avatar, &info.CreatedDate)
	//fmt.Println(info)
	return
}

func (d *dao) RawUserInfoID(ctx context.Context, id int64) (info *model.UserInfo, err error) {
	info = new(model.UserInfo)
	err = d.db.QueryRow(ctx, _selUserInfoID, id).Scan(&info.UserID, &info.Name, &info.Tel, &info.Mail, &info.Gender, &info.Avatar, &info.Description, &info.CreatedDate)
	if err != nil {
		log.Error("d.RawInfo.Query error by id = %d, (%v)", id, err)
		return
	}
	return
}

func (d *dao) RawUserInfoName(ctx context.Context, name string) (info *model.UserInfo, err error) {
	info = new(model.UserInfo)
	err = d.db.QueryRow(ctx, _selUserInfoName, name).Scan(&info.UserID, &info.Name, &info.Tel, &info.Mail, &info.Gender, &info.Avatar, &info.Description, &info.CreatedDate)
	if err != nil {
		log.Error("d.RawInfo.Query error by name = %s, (%v)", name, err)
		return
	}
	return
}