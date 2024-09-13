package fncache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
)

type NormalFn[I, O any] func(ctx context.Context, in I) (out O, err error)

// CacheFn 包装函数，增加缓存功能
//
//	getkey 获取缓存key
//	when   当返回true时，记录缓存
func CacheFn[I, O any](rdb redis.UniversalClient, fn NormalFn[I, O], getkey func(I) string, when func(I, O) bool, expire time.Duration) NormalFn[I, O] {
	return func(ctx context.Context, in I) (O, error) {

		key := getkey(in)
		if s, err2 := rdb.Get(ctx, key).Result(); err2 == nil && len(s) > 0 {
			var out O
			if err3 := json.Unmarshal([]byte(s), &out); err3 == nil {
				return out, nil
			}
		}

		out, err := fn(ctx, in)
		if err != nil {
			return out, err
		}

		if when(in, out) {
			if bs, err2 := json.Marshal(out); err2 == nil {
				rdb.Set(ctx, key, string(bs), expire)
			}
		}

		return out, err
	}
}
