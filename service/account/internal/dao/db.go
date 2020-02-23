package dao

import (
	"context"
	"fmt"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/database/sql"
	"github.com/bilibili/kratos/pkg/log"
	"account/internal/model"
)

const (
	_selUserInfoID = "SELECT uid,name,tel,mail,gender,avatar,description,created FROM user WHERE uid=? "
	_selUserInfoTel = "SELECT uid,name,tel,mail,gender,avatar,description,created FROM user WHERE tel=? "
	_selUserInfoName = "SELECT uid,name,tel,mail,gender,avatar,description,created FROM user WHERE name=?"

	_searchUserName = "SELECT * from user where name like ?" // %str%

)

func NewDB() (db *sql.DB, err error) {
	var cfg struct {
		Client *sql.Config
	}
	if err = paladin.Get("db.toml").UnmarshalTOML(&cfg); err != nil {
		return
	}
	db = sql.NewMySQL(cfg.Client)
	err = db.Ping(context.Background())
	if err != nil {
		panic(err)
	}
	return
}

func (d *Dao) RawUserInfoID(ctx context.Context, id int64) (info *model.UserInfo, err error) {
	info = new(model.UserInfo)
	err = d.db.QueryRow(ctx, _selUserInfoID, id).Scan(&info.UserID, &info.Name, &info.Tel, &info.Mail, &info.Gender, &info.Avatar, &info.Description, &info.CreatedDate)
	if err != nil && err != sql.ErrNoRows {
		log.Error("d.RawInfo.Query error by id = %d, (%v)", id, err)
		return
	}
	return
}

func (d *Dao) UserInfoTel(ctx context.Context, tel string) (info *model.UserInfo, err error) {
	info = new(model.UserInfo)
	err = d.db.QueryRow(ctx, _selUserInfoTel, tel).Scan(&info.UserID, &info.Name, &info.Tel, &info.Mail, &info.Gender, &info.Avatar, &info.Description, &info.CreatedDate)
	if err != nil {
		log.Error("d.RawInfo.Query error by tel = %d, (%v)", tel, err)
		return
	}
	return
}

func (d *Dao) RawUserInfoName(ctx context.Context, name string) (info *model.UserInfo, err error) {
	info = new(model.UserInfo)
	err = d.db.QueryRow(ctx, _selUserInfoName, name).Scan(&info.UserID, &info.Name, &info.Tel, &info.Mail, &info.Gender, &info.Avatar, &info.Description, &info.CreatedDate)
	if err != nil && err != sql.ErrNoRows {
		log.Error("d.RawInfo.Query error by name = %s, (%v)", name, err)
		return
	}
	return
}

func (d *Dao) SearchUserInfoByName(c context.Context, name string) (infos []model.UserInfo, err error) {
	name = fmt.Sprintf("%%%s%%", name)
	rows, err := d.db.Query(c, _searchUserName, name)
	infos = []model.UserInfo{}
	if err != nil && err != sql.ErrNoRows {
		log.Error("query  error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var info model.UserInfo = model.UserInfo{}
		err = rows.Scan(&info.UserID, &info.Name, &info.Tel, &info.Mail, &info.Gender, &info.Avatar, &info.Description, &info.CreatedDate)
		if err != nil {
			log.Error("scan demo log error(%v)", err)
			return
		}
		infos = append(infos, info)
	}
	d.cache.Do(c, func(ctx context.Context) {
		for _, info := range infos {
			d.AddCacheUserInfo(c, info.UserID, &info)
		}
	})
	return
}