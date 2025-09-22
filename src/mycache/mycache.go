package mycache

import (
	"context"
	"errors"
	"log"
	"orders/src/metrics"
	"os"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/extra/redisprometheus/v9"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/sdk/trace"
)

type RedisService struct {
	client *cache.Cache
	rdb    *redis.Client
	ttl    time.Duration
	m      *metrics.Metrics
}

func NewRedis(tp *trace.TracerProvider, reg prometheus.Registerer, m *metrics.Metrics, ttl time.Duration) *RedisService {
	redisAddr := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	collector := redisprometheus.NewCollector("redis", "client", rdb)

	reg.MustRegister(collector)

	redisotel.WithTracerProvider(tp)

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		log.Printf("Error in InstrumentTracing: %v\n", err)
	}

	if err := redisotel.InstrumentMetrics(rdb); err != nil {
		log.Printf("Error in InstrumentMetrics: %v\n", err)
	}

	mycache := cache.New(&cache.Options{
		Redis:      rdb,
		LocalCache: cache.NewTinyLFU(1000, ttl),
	})

	return &RedisService{client: mycache, rdb: rdb, ttl: time.Hour, m: m}
}

func (r *RedisService) Get(ctx context.Context, key string, value interface{}) error {
	err := r.client.Get(ctx, key, value)

	if err == nil {
		r.m.CacheHits.Inc()
	} else if errors.Is(err, cache.ErrCacheMiss) {
		r.m.CacheMisses.Inc()
	}

	return err
}

func (r *RedisService) Set(ctx context.Context, key string, value interface{}) error {
	err := r.client.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: &value,
		TTL:   r.ttl,
	})

	return err
}

func (r *RedisService) Close() error {
	return r.rdb.Close()
}
