package bfcache_test

import (
	"context"
	"encoding/binary"
	"math/rand"
	"testing"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/allegro/bigcache/v3"
	"github.com/stretchr/testify/require"
)

const (
	cacheSize       = 9 * 1024 * 1024 * 1024 // 9 GB
	cacheEntrySize  = 800                    // 800 B
	cacheEntryCount = 8_000_000
)

func BenchmarkBigCacheGet(b *testing.B) {
	ctx := context.Background()
	cache, err := bigcache.New(ctx, bigcache.Config{
		Shards:           1024,
		HardMaxCacheSize: cacheSize,
		MaxEntrySize:     cacheEntrySize,
		StatsEnabled:     true,
	})
	require.NoError(b, err)
	benchmarkGenericGet(b, &bigcacheWrapper{cache})
}

func BenchmarkFastCacheGet(b *testing.B) {
	cache := fastcache.New(cacheSize)
	benchmarkGenericGet(b, &fastcacheWrapper{cache})
}

func benchmarkGenericGet(b *testing.B, cache GenericCache) {
	rng := rand.New(rand.NewSource(0))

	key := make([]byte, binary.MaxVarintLen64)
	value := make([]byte, cacheEntrySize)

	for i := uint64(0); i < cacheEntryCount; i++ {
		binary.NativeEndian.PutUint64(key, i)
		rng.Read(value[:0])
		cache.Set(key, value)
	}

	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		key := make([]byte, binary.MaxVarintLen64)
		i := uint64(0)
		for p.Next() {
			binary.NativeEndian.PutUint64(key, i)
			i = (i + 400) % cacheEntryCount

			cache.Get(key)
		}
	})
}
