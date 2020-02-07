package dao

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"kratos-demo/internal/model"
	"github.com/bilibili/kratos/pkg/cache/memcache"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/log"
)

//go:generate kratos tool genmc
type _mc interface {
	// mc: -key=keyInfo -type=get
	CacheUserInfo(c context.Context, key interface{}) (*model.UserInfo, error)
	// mc: -key=keyInfo -expire=d.demoExpire
	AddCacheUserInfo(c context.Context, key interface{}, art *model.UserInfo) (err error)
	// mc: -key=keyInfo
	DeleteUserInfoCache(c context.Context, key interface{}) (err error)
}

func NewMC() (mc *memcache.Memcache, err error) {
	var cfg struct {
		Client *memcache.Config
	}
	if err = paladin.Get("memcache.toml").UnmarshalTOML(&cfg); err != nil {
		return
	}
	mc =  memcache.New(cfg.Client)
	return
}

func (d *dao) PingMC(ctx context.Context) (err error) {
	if err = d.mc.Set(ctx, &memcache.Item{Key: "ping", Value: []byte("pong"), Expiration: 0}); err != nil {
		log.Error("conn.Set(PING) error(%v)", err)
	}
	return
}

// FIXME: what if the id is the same to name?
func keyInfo(key interface{}) string {
	if strings.Compare(reflect.TypeOf(key).String(), "string") == 0 {
		return fmt.Sprintf("user_%s", key)
	} else {
		return fmt.Sprintf("user_%d", key)
	}
}
