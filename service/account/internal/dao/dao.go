package dao

import (
	"context"
	"time"

	"github.com/bilibili/kratos/pkg/cache/memcache"
	"github.com/bilibili/kratos/pkg/cache/redis"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/database/sql"
	"github.com/bilibili/kratos/pkg/sync/pipeline/fanout"
	xtime "github.com/bilibili/kratos/pkg/time"
	"account/internal/model"
)

//go:generate kratos tool genbts
// Dao dao interface
type bts interface {
	Close()
	Ping(ctx context.Context) (err error)
	// bts: -nullcache=&model.UserInfo{UserID:-1} -check_null_code=$!=nil&&$.UserID==-1
	UserInfoID(c context.Context, id int64) (*model.UserInfo, error)
	// bts: -nullcache=&model.UserInfo{UserID:-1} -check_null_code=$!=nil&&$.UserID==-1
	UserInfoName(c context.Context, name string) (*model.UserInfo, error)
}

// dao dao.
type Dao struct {
	db          *sql.DB
	redis       *redis.Redis
	mc          *memcache.Memcache
	cache 		*fanout.Fanout
	demoExpire int32
}

// New new a dao and return.
func New(r *redis.Redis, mc *memcache.Memcache, db *sql.DB) (d *Dao, err error) {
	var cfg struct{
		DemoExpire xtime.Duration
	}
	if err = paladin.Get("application.toml").UnmarshalTOML(&cfg); err != nil {
		return
	}
	d = &Dao{
		db: db,
		redis: r,
		mc: mc,
		cache: fanout.New("cache"),
		demoExpire: int32(time.Duration(cfg.DemoExpire) / time.Second),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.mc.Close()
	d.redis.Close()
	d.db.Close()
	d.cache.Close()
}

// Ping ping the resource.
func (d *Dao) Ping(ctx context.Context) (err error) {
	return nil
}
