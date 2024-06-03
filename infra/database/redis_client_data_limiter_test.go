package database

import (
	"context"
	"testing"
	"time"

	"github.com/jpodlasnisky/ratelimiter/config"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func setup() (*RedisDataLimiter, func()) {
	// Inicie um servidor miniredis.
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	// Crie um cliente Redis que se conecta ao miniredis.
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// Crie um RedisDataLimiter com o cliente Redis.
	limiter := NewRedisDataLimiter(client)

	return limiter, func() {
		mr.Close()
	}
}

func TestNewRedisClient(t *testing.T) {
	cfg := &config.Config{RedisURL: "localhost:6379"}
	client := NewRedisClient(cfg)

	assert.NotNil(t, client)
	assert.Equal(t, cfg.RedisURL, client.Options().Addr)
	assert.Equal(t, "", client.Options().Password)
	assert.Equal(t, 0, client.Options().DB)
}

func TestZAdd(t *testing.T) {
	limiter, teardown := setup()
	defer teardown()

	ctx := context.Background()

	members := []*redis.Z{{Score: 1, Member: "member1"}}
	added, err := limiter.ZAdd(ctx, "key1", members...)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), added)
}

func TestZCard(t *testing.T) {
	limiter, teardown := setup()
	defer teardown()

	ctx := context.Background()

	members := []*redis.Z{{Score: 1, Member: "member1"}}
	_, err := limiter.ZAdd(ctx, "key1", members...)
	assert.NoError(t, err)

	count, err := limiter.ZCard(ctx, "key1")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestExists(t *testing.T) {
	limiter, teardown := setup()
	defer teardown()

	ctx := context.Background()

	err := limiter.Set(ctx, "key1", "value1", time.Minute)
	assert.NoError(t, err)

	exists, err := limiter.Exists(ctx, "key1")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), exists)
}

func TestSet(t *testing.T) {
	limiter, teardown := setup()
	defer teardown()

	ctx := context.Background()

	err := limiter.Set(ctx, "key2", "value2", time.Minute)
	assert.NoError(t, err)
}

func TestGet(t *testing.T) {
	limiter, teardown := setup()
	defer teardown()

	ctx := context.Background()

	err := limiter.Set(ctx, "key2", "value2", time.Minute)
	assert.NoError(t, err)

	value, err := limiter.Get(ctx, "key2")
	assert.NoError(t, err)
	assert.Equal(t, "value2", value)
}

func TestSetEX(t *testing.T) {
	limiter, teardown := setup()
	defer teardown()

	ctx := context.Background()

	err := limiter.SetEX(ctx, "key3", "value3", time.Minute)
	assert.NoError(t, err)
}

func TestZRemRangeByScore(t *testing.T) {
	limiter, teardown := setup()
	defer teardown()

	ctx := context.Background()

	members := []*redis.Z{{Score: 1, Member: "member1"}}
	_, err := limiter.ZAdd(ctx, "key1", members...)
	assert.NoError(t, err)

	removed, err := limiter.ZRemRangeByScore(ctx, "key1", "0", "1")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), removed)
}
