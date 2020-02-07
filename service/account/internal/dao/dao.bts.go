package dao

import (
	"context"

	"github.com/bilibili/kratos/pkg/cache"
	"kratos-demo/internal/model"
)

// UserInfoID get data from cache if miss will call source method, then add to cache.
func (d *dao) UserInfoID(c context.Context, id int64) (res *model.UserInfo, err error) {
	addCache := true
	res, err = d.CacheUserInfo(c, id)
	if err != nil {
		addCache = false
		err = nil
	}
	defer func() {
		if res != nil && res.UserID == -1 {
			res = nil
		}
	}()
	if res != nil {
		cache.MetricHits.Inc("bts:UserInfoID")
		return
	}
	cache.MetricMisses.Inc("bts:UserInfoID")
	res, err = d.RawUserInfoID(c, id)
	if err != nil {
		return
	}
	miss := res
	if miss == nil {
		miss = &model.UserInfo{UserID: -1}
	}
	if !addCache {
		return
	}
	d.cache.Do(c, func(c context.Context) {
		d.AddCacheUserInfo(c, id, miss)
	})
	return
}

// UserInfoName get data from cache if miss will call source method, then add to cache.
func (d *dao) UserInfoName(c context.Context, name string) (res *model.UserInfo, err error) {
	addCache := true
	res, err = d.CacheUserInfo(c, name)
	if err != nil {
		addCache = false
		err = nil
	}
	defer func() {
		if res != nil && res.UserID == -1 {
			res = nil
		}
	}()
	if res != nil {
		cache.MetricHits.Inc("bts:UserInfoName")
		return
	}
	cache.MetricMisses.Inc("bts:UserInfoName")
	res, err = d.RawUserInfoName(c, name)
	if err != nil {
		return
	}
	miss := res
	if miss == nil {
		miss = &model.UserInfo{UserID: -1}
	}
	if !addCache {
		return
	}
	d.cache.Do(c, func(c context.Context) {
		d.AddCacheUserInfo(c, name, miss)
	})
	return
}
