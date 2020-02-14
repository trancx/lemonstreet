package dao

import (
	"comment/api/cmtapi"
	"context"
	"github.com/bilibili/kratos/pkg/naming/discovery"
	"github.com/bilibili/kratos/pkg/net/rpc/warden/resolver"
	"time"

	"github.com/bilibili/kratos/pkg/cache/memcache"
	"github.com/bilibili/kratos/pkg/cache/redis"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/database/sql"
	"github.com/bilibili/kratos/pkg/sync/pipeline/fanout"
	xtime "github.com/bilibili/kratos/pkg/time"
)

//go:generate kratos tool genbts
// Dao dao interface
type Dao interface {
	Close()
	Ping(ctx context.Context) (err error)
	//// bts: -nullcache=&model.Article{ID:-1} -check_null_code=$!=nil&&$.ID==-1
	//Article(c context.Context, id int64) (*model.Article, error)
	PostComment(ctx context.Context, comment *cmtapi.Comment) (err error)
	SearchCommentsByUId(c context.Context, uid int64) (cmmts []*cmtapi.Comment, err error)
	SearchCommentsByAId(c context.Context, aid int64) (cmmts []*cmtapi.Comment, err error)
}

func init(){
	// NOTE: 注意这段代码，表示要使用discovery进行服务发现
	// NOTE: 还需注意的是，resolver.Register是全局生效的，所以建议该代码放在进程初始化的时候执行
	// NOTE: ！！！切记不要在一个进程内进行多个不同中间件的Register！！！
	// NOTE: 在启动应用时，可以通过flag(-discovery.nodes) 或者 环境配置(DISCOVERY_NODES)指定discovery节点
	resolver.Register(discovery.Builder())
}

// comment  should have an id, represents the rank of the article
// dao dao.
type dao struct {
	db          *sql.DB
	redis       *redis.Redis
	mc          *memcache.Memcache
	cache *fanout.Fanout
	demoExpire int32
}

// New new a dao and return.
func New(r *redis.Redis, mc *memcache.Memcache, db *sql.DB) (d Dao, err error) {
	var cfg struct{
		DemoExpire xtime.Duration
	}
	if err = paladin.Get("application.toml").UnmarshalTOML(&cfg); err != nil {
		return
	}
	d = &dao{
		db: db,
		redis: r,
		mc: mc,
		cache: fanout.New("cache"),
		demoExpire: int32(time.Duration(cfg.DemoExpire) / time.Second),
	}
	return
}

// Close close the resource.
func (d *dao) Close() {
	d.mc.Close()
	d.redis.Close()
	d.db.Close()
	d.cache.Close()
}

// Ping ping the resource.
func (d *dao) Ping(ctx context.Context) (err error) {
	return nil
}
