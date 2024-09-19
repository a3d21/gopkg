package rdbutil

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func SaddAndReturn(ctx context.Context, rdb redis.UniversalClient, key string, member string) (add, total int64, err error) {
	script := `-- luaSaddAndReturnSCard
local ret = {}
ret[1] = redis.call('SADD', KEYS[1], ARGV[1])
ret[2] = redis.call('SCARD', KEYS[1])
return ret`

	ret, err := rdb.Eval(ctx, script, []string{key}, []string{member}).Result()
	if err != nil {
		return 0, 0, err
	}
	if vs, ok := ret.([]any); ok && len(vs) == 2 {
		add, _ = vs[0].(int64)
		total, _ = vs[1].(int64)
	}
	return add, total, nil
}
