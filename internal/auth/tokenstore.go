package tokenstore

import (
    "context"
    "time"

    "github.com/redis/go-redis/v9"
)

type RedisStore struct {
    rdb *redis.Client
}

func NewRedisStore(addr, pass string, db int) *RedisStore {
    r := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: pass,
        DB:       db,
    })
    return &RedisStore{rdb: r}
}

// Save a token with TTL
func (s *RedisStore) Save(ctx context.Context, token string, ttl time.Duration) error {
    return s.rdb.Set(ctx, token, "1", ttl).Err()
}

// Exists returns true if token is still in store
func (s *RedisStore) Exists(ctx context.Context, token string) (bool, error) {
    n, err := s.rdb.Exists(ctx, token).Result()
    return n == 1, err
}

// Delete revokes a token immediately
func (s *RedisStore) Delete(ctx context.Context, token string) error {
    return s.rdb.Del(ctx, token).Err()
}
