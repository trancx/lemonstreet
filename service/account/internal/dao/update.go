package dao

import (
	"context"
	"github.com/bilibili/kratos/pkg/cache"
	"github.com/bilibili/kratos/pkg/log"

	"kratos-demo/internal/model"
)

const (
	_setUserInfo = "INSERT user (`uid`, `name`, `tel`, `mail`, `gender`, `avatar`, `description`, `created`) VALUES (?,?,?,?,?,?,?,?)"
	_updateAvatar = "UPDATE user SET avatar=? WHERE (uid=?)"
	_updateDesc = "UPDATE user SET description=? WHERE (uid=?)"
	_updateGender = "UPDATE user SET gender=? WHERE (uid=?)"
	_updateMail = "UPDATE user SET mail=? WHERE (uid=?)"
)

func (d *dao) SetUserInfo(c context.Context, info *model.UserInfo) (err error) {
	var (
		id int64
	)
	id = info.UserID
	_, err = d.db.Exec(c, _setUserInfo, info.UserID, info.Name, info.Tel, info.Mail, info.Gender, info.Avatar, info.Description, info.CreatedDate)
	
	if err != nil {
		log.Error("db.SetUserInfo.Exec(%s) error(%v)", _setUserInfo, err)
		return
	}
	d.cache.Do(c, func(ctx context.Context) {
		d.AddCacheUserInfo(ctx, id, info)
	})
	return
}

func (d *dao) UpdateAvatar(c context.Context, id int64, avatar string) (err error) {
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
		log.Error("db.UpdateAvatar.Exec(%s) error(%v)", _updateAvatar, err)
		return
	}
	log.Info("db.UpdateAvatar.Exec(%s) success", _updateAvatar)
	return 
}

func (d *dao) UpdateDesc(c context.Context, id int64, desc string) (err error) {
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
		log.Error("db.UpdateDesc.Exec(%s) error(%v)", _updateDesc, err)
		return
	}
	log.Info("db.UpdateDesc.Exec(%s) success", _updateDesc)
	return 
}

func (d *dao) UpdateGender(c context.Context, id int64, gender string) (err error) {
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
		log.Error("db.UpdateGender.Exec(%s) error(%v)", _updateGender, err)
		return
	}
	log.Info("db.UpdateGender.Exec(%s) success", _updateGender)
	return 
}

func (d *dao) UpdateMail(c context.Context, id int64, mail string) (err error) {
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
		log.Error("db.UpdateMail.Exec(%s) error(%v)", _updateMail, err)
		return
	}
	log.Info("db.UpdateMail.Exec(%s) success", _updateMail)
	return 
}
