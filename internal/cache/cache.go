package cache

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	Def *gcache.Cache
)

func init() {
	ctx := gctx.New()
	cacheAdapter, err := g.Config().Get(ctx, "cache.adapter")
	if err != nil {
		panic(err)
	}
	Def = gcache.New()
	if cacheAdapter.String() == "redis" {
		host, err := g.Config().Get(ctx, "redis.cache.address")
		db, err := g.Config().Get(ctx, "redis.cache.db")
		pass, err := g.Config().Get(ctx, "redis.cache.pass")
		if err != nil {
			panic(err)
		}
		redis, err := gredis.New(&gredis.Config{
			Address: host.String(),
			Db:      db.Int(),
			Pass:    pass.String(),
		})
		if err != nil {
			panic(err)
		}
		Def.SetAdapter(gcache.NewAdapterRedis(redis))
	} else {
		Def.SetAdapter(gcache.New())
	}

}
