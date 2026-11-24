package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// GetOrSet mengambil data dari Redis jika ada (hit), atau memanggil fetchFn jika miss lalu menyimpannya ke Redis.
// Jika Redis nil atau error, fungsi akan graceful fallback ke fetchFn tanpa menggagalkan flow request.
func GetOrSet[T any](
	ctx context.Context,
	rdb *redis.Client,
	key string,
	ttl time.Duration,
	fetchFn func() (T, error),
) (T, error) {
	var zero T

	// Ambil dari Redis jika client tersedia
	if rdb != nil {
		if val, err := rdb.Get(ctx, key).Result(); err == nil {
			var result T
			if err := json.Unmarshal([]byte(val), &result); err == nil {
				return result, nil // Cache Hit
			}
		}
	}

	// 2. Cache Miss / Redis tidak tersedia -> panggil DB atau kalkulasi via fetchFn
	result, err := fetchFn()
	if err != nil {
		return zero, err
	}

	// Simpan hasil ke Redis
	if rdb != nil {
		if bytes, err := json.Marshal(result); err == nil {
			_ = rdb.Set(ctx, key, bytes, ttl).Err()
		}
	}

	return result, nil
}

// Invalidate menghapus satu atau beberapa key spesifik dari Redis.
func Invalidate(ctx context.Context, rdb *redis.Client, keys ...string) {
	if rdb != nil && len(keys) > 0 {
		_ = rdb.Del(ctx, keys...).Err()
	}
}

// InvalidateByPattern menghapus semua key yang cocok dengan pola (contoh: "master:models:brand:*").
func InvalidateByPattern(ctx context.Context, rdb *redis.Client, pattern string) {
	if rdb != nil {
		if keys, err := rdb.Keys(ctx, pattern).Result(); err == nil && len(keys) > 0 {
			_ = rdb.Del(ctx, keys...).Err()
		}
	}
}
