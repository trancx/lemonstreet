package dao

import (
	"context"
	"database/sql"
	"github.com/bilibili/kratos/pkg/cache"
	"github.com/bilibili/kratos/pkg/log"

	"account/internal/model"
)

const (
	_setUserInfo = "INSERT user (`name`, `tel`, `created`) VALUES (?,?,?)"
	_updateAvatar = "UPDATE user SET avatar=? WHERE (uid=?)"
	_updateDesc = "UPDATE user SET description=? WHERE (uid=?)"
	_updateGender = "UPDATE user SET gender=? WHERE (uid=?)"
	_updateMail = "UPDATE user SET mail=? WHERE (uid=?)"
	_updateName = "UPDATE user SET name=? WHERE (uid=?)"
)

// GENARATE: uid, name conflic?
func (d *Dao) InsertUserInfo(c context.Context, info *model.UserInfo) (id int64, err error) {
	var (
		rst sql.Result
	)
	rst, err = d.db.Exec(c, _setUserInfo, info.Name, info.Tel, info.CreatedDate)
	if err != nil {
		log.Error("db.InsertUserInfo.Exec(%s) error(%v)", _setUserInfo, err)
		return
	}
	id, err = rst.LastInsertId()
	if err != nil {
		log.Error("db.InsertUserInfo.Exec(%s) error(%v)", _setUserInfo, err)
		return
	}
	d.cache.Do(c, func(ctx context.Context) {
		d.AddCacheUserInfo(ctx, id, info)
	})
	return
}

func (d *Dao) UpdateName(c context.Context, id int64, name string) (err error) {
	var (
		info 	*model.UserInfo
	)
	info, err = d.CacheUserInfo(c, id)
	if info != nil {
		cache.MetricHits.Inc("bts:UserInfoID")
		err = d.DeleteUserInfoCache(c, id)
		if err != nil {
			info.Name = name
			d.cache.Do(c, func(ctx context.Context) {
				d.AddCacheUserInfo(ctx, id, info)
			})

		}
	}
	// access avatar FIXME: result should be handled
	_, err = d.db.Exec(c, _updateName, name, id)
	if err != nil {
		log.Error("db.UpdateName.Exec error(%v)", err)
		return
	}
	log.Info("db.UpdateName.Exec success")
	return
}

func (d *Dao) UpdateAvatar(c context.Context, id int64, avatar string) (err error) {
	var (
		info 	*model.UserInfo
	)
	info, err = d.CacheUserInfo(c, id)
	if info != nil {
		cache.MetricHits.Inc("bts:UserInfoID")
		err = d.DeleteUserInfoCache(c, id)
		if err != nil {
			info.Avatar = avatar
			d.cache.Do(c, func(ctx context.Context) {
				d.AddCacheUserInfo(ctx, id, info)
			})

		}
	}
	// access avatar FIXME: result should be handled
	_, err = d.db.Exec(c, _updateAvatar, avatar, id)
	if err != nil {
		log.Error("db.UpdateAvatar.Exec error(%v)", err)
		return
	}
	log.Info("db.UpdateAvatar.Exec success")
	return 
}

func (d *Dao) UpdateDesc(c context.Context, id int64, desc string) (err error) {
	var (
		info 	*model.UserInfo
	)
	info, err = d.CacheUserInfo(c, id)
	if info != nil {
		cache.MetricHits.Inc("bts:UserInfoID")
		err = d.DeleteUserInfoCache(c, id)
		if err != nil {
			info.Description = desc
			d.cache.Do(c, func(ctx context.Context) {
				d.AddCacheUserInfo(ctx, id, info)
			})
		}
	}
	// FIXME: result should be handled
	_, err = d.db.Exec(c, _updateDesc, desc, id)
	if err != nil {
		log.Error("db.UpdateDesc.Exec error(%v)", err)
		return
	}
	log.Info("db.UpdateDesc.Exec success")
	return 
}

func (d *Dao) UpdateGender(c context.Context, id int64, gender string) (err error) {
	var (
		info 	*model.UserInfo
	)
	info, err = d.CacheUserInfo(c, id)
	if info != nil {
		cache.MetricHits.Inc("bts:UserInfoID")
		err = d.DeleteUserInfoCache(c, id)
		if err != nil {
			info.Gender = gender
			d.cache.Do(c, func(ctx context.Context) {
				d.AddCacheUserInfo(ctx, id, info)
			})
		}
	}
	// FIXME: result should be handled
	_, err = d.db.Exec(c, _updateGender, gender, id)
	if err != nil {
		log.Error("db.UpdateGender.Exec error(%v)", err)
		return
	}
	log.Info("db.UpdateGender.Exec success")
	return 
}

func (d *Dao) UpdateMail(c context.Context, id int64, mail string) (err error) {
	var (
		info 	*model.UserInfo
	)
	info, err = d.CacheUserInfo(c, id)
	if info != nil {
		cache.MetricHits.Inc("bts:UserInfoID")
		err = d.DeleteUserInfoCache(c, id)
		if err != nil {
			info.Mail = mail
			err = d.AddCacheUserInfo(c, id, info)
		}
	}
	// FIXME: result should be handled
	_, err = d.db.Exec(c, _updateMail, mail, id)
	if err != nil {
		log.Error("db.UpdateMail.Exec error(%v)", err)
		return
	}
	log.Info("db.UpdateMail.Exec success")
	return 
}
