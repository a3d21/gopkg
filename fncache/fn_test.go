package fncache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	fnCallCount int = 0
)

func TestCacheFn(t *testing.T) {
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{"localhost:6379"},
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	getkey := func(s *SomeReq) string {
		return fmt.Sprintf("key:%v", s.ID)
	}
	when := func(s *SomeReq, s2 *SomeResp) bool {
		return s.ID < 10
	}
	cachedFn := CacheFn(rdb, someFn1, getkey, when, 100*time.Millisecond)

	ctx := context.Background()
	_, err := cachedFn(ctx, &SomeReq{ID: 3})
	assert.Nil(t, err)
	assert.Equal(t, 1, fnCallCount)

	// use cache
	_, err = cachedFn(ctx, &SomeReq{ID: 3})
	assert.Nil(t, err)
	assert.Equal(t, 1, fnCallCount)

	// id=13 no cache
	_, err = cachedFn(ctx, &SomeReq{ID: 13})
	assert.Nil(t, err)
	assert.Equal(t, 2, fnCallCount)

	_, err = cachedFn(ctx, &SomeReq{ID: 13})
	assert.Nil(t, err)
	assert.Equal(t, 3, fnCallCount)
}

type SomeReq struct {
	ID int64
}

type SomeResp struct {
	Name string
}

func someFn1(ctx context.Context, req *SomeReq) (*SomeResp, error) {
	fnCallCount += 1
	return &SomeResp{Name: fmt.Sprint(req.ID)}, nil
}
