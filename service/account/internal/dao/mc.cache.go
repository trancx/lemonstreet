// Code generated by kratos tool genmc. DO NOT EDIT.

/*
  Package dao is a generated mc cache package.
  It is generated from:
  type _mc interface {
		// mc: -key=keyInfo -type=get
		CacheUserInfo(c context.Context, key interface{}) (*model.UserInfo, error)
		// mc: -key=keyInfo -expire=d.demoExpire
		AddCacheUserInfo(c context.Context, key interface{}, art *model.UserInfo) (err error)
		// mc: -key=keyInfo
		DeleteUserInfoCache(c context.Context, key interface{}) (err error)
	}
*/

package dao

import (
	"context"
	"fmt"

	"github.com/bilibili/kratos/pkg/cache/memcache"
	"github.com/bilibili/kratos/pkg/log"
	"account/internal/model"
)

var _ _mc

// CacheUserInfo get data from mc
func (d *Dao) CacheUserInfo(c context.Context, id interface{}) (res *model.UserInfo, err error) {
	key := keyInfo(id)
	res = &model.UserInfo{}
	if err = d.mc.Get(c, key).Scan(res); err != nil {
		res = nil
		if err == memcache.ErrNotFound {
			err = nil
		}
	}
	if err != nil {
		log.Errorv(c, log.KV("CacheUserInfo", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}

// AddCacheUserInfo Set data to mc
func (d *Dao) AddCacheUserInfo(c context.Context, id interface{}, val *model.UserInfo) (err error) {
	if val == nil {
		return
	}
	key := keyInfo(id)
	item := &memcache.Item{Key: key, Object: val, Expiration: d.demoExpire, Flags: memcache.FlagJSON}
	if err = d.mc.Set(c, item); err != nil {
		log.Errorv(c, log.KV("AddCacheUserInfo", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}

// DeleteUserInfoCache delete data from mc
func (d *Dao) DeleteUserInfoCache(c context.Context, id interface{}) (err error) {
	key := keyInfo(id)
	if err = d.mc.Delete(c, key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Errorv(c, log.KV("DeleteUserInfoCache", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}
