package mycache

import (
	"context"
	"orders/src/metrics"
	"testing"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/go-redis/redismock/v9"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/redis/go-redis/v9"

	"github.com/stretchr/testify/require"
)

func newTestMetrics() *metrics.Metrics {
	reg := prometheus.NewRegistry()
	return metrics.New(reg, reg)
}

func newMockRedis() (CacheService, redismock.ClientMock, *metrics.Metrics) {

	TTL := 1000 * time.Second

	m := newTestMetrics()

	db, mock := redismock.NewClientMock()

	mockcache := cache.New(&cache.Options{
		Redis:      db,
		LocalCache: nil,
	})

	return &redisService{
		rdb:    db,
		ttl:    TTL,
		client: mockcache,
		m:      m,
	}, mock, m
}

func TestCacheGet(t *testing.T) {
	ctx := context.Background()
	service, mock, m := newMockRedis()

	mock.ExpectGet("no-data").SetErr(cache.ErrCacheMiss)
	mock.ExpectGet("has-data").SetVal("some-data")

	// Проверяем ошибку для несуществующего ключа - должна упасть ошибка и инкремент CacheMisses

	var empty interface{}
	var expected interface{}

	err := service.Get(ctx, "no-data", &empty)

	missCouner := testutil.ToFloat64(m.CacheMisses)
	require.Equal(t, float64(1), missCouner)

	require.Equal(t, expected, empty)
	require.EqualError(t, cache.ErrCacheMiss, err.Error())

	var hasData string

	err = service.Get(ctx, "has-data", &hasData)

	require.NoError(t, err)

	hitsCounter := testutil.ToFloat64(m.CacheHits)

	require.Equal(t, "some-data", hasData)
	require.Equal(t, float64(1), hitsCounter)

}

func TestClosed(t *testing.T) {
	ctx := context.Background()
	service, mock, metrics := newMockRedis()

	mock.ExpectGet("closed").SetErr(redis.ErrClosed)

	var empty interface{}

	err := service.Get(ctx, "closed", empty)

	hitsConter := testutil.ToFloat64(metrics.CacheHits)
	missConter := testutil.ToFloat64(metrics.CacheMisses)

	require.Equal(t, float64(0), hitsConter)
	require.Equal(t, float64(0), missConter)
	require.EqualError(t, redis.ErrClosed, err.Error())

}
