package dao

import (
	"context"

	"kratos-demo/internal/model"
)

const (
	_setUserInfo = "INSERT"
	_updateAvatar = "UPDATE"
	_updateDesc = "UPDATE"
	_updateGender = "UPDATE"
	_updateMail = "UPDATE"
)

func (d *dao) SetUserInfo(c context.Context, info *model.UserInfo) error {
	panic("implement me")
}

func (d *dao) UpadateAvatar(c context.Context, id int64, avatar string) error {
	var (
		info 	*model.UserInfo
		err 	error
	)
	info, err = d.CacheUserInfo(c, id)
	if info != nil {
		err = d.DeleteUserInfoCache(c, id)
		if err != nil {
			info.Avatar = avatar
			err = d.AddCacheUserInfo(c, id, info)
		}
	}
	// access avatar
	d.db.Prepare()
}

func (d *dao) UpadateDesc(c context.Context, id int64, desc string) error {
	var (
		info 	*model.UserInfo
		err 	error
	)
	info, err = d.CacheUserInfo(c, id)
	if info != nil {
		err = d.DeleteUserInfoCache(c, id)
		if err != nil {
			info.Avatar = avatar
			err = d.AddCacheUserInfo(c, id, info)
		}
	}
	// access avatar
	d.db.Prepare()
}

func (d *dao) UpadateGender(c context.Context, id int64, gender string) error {
	var (
		info 	*model.UserInfo
		err 	error
	)
	info, err = d.CacheUserInfo(c, id)
	if info != nil {
		err = d.DeleteUserInfoCache(c, id)
		if err != nil {
			info.Avatar = avatar
			err = d.AddCacheUserInfo(c, id, info)
		}
	}
	// access avatar
	d.db.Prepare()
}

func (d *dao) UpadateMail(c context.Context, id int64, mail string) error {
	var (
		info 	*model.UserInfo
		err 	error
	)
	info, err = d.CacheUserInfo(c, id)
	if info != nil {
		err = d.DeleteUserInfoCache(c, id)
		if err != nil {
			info.Avatar = avatar
			err = d.AddCacheUserInfo(c, id, info)
		}
	}
	// access avatar
	d.db.Prepare()
}
