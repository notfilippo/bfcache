package bfcache_test

import (
	"context"
	"encoding/binary"
	"math/rand"
	"testing"
	"time"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/allegro/bigcache/v3"
	"github.com/stretchr/testify/require"
)

const (
	cacheSize       = 9 * 1024 * 1024 * 1024 // 9 GB
	cacheEntrySize  = 800                    // 800 B
	cacheEntryCount = 8_000_000
)

func newBigCache(tb testing.TB) GenericCache {
	ctx := context.Background()
	cache, err := bigcache.New(ctx, bigcache.Config{
		Shards:           1024,
		LifeWindow:       1 * time.Hour,
		HardMaxCacheSize: cacheSize,
		MaxEntrySize:     cacheEntrySize,
		StatsEnabled:     true,
	})
	require.NoError(tb, err)
	return &bigcacheWrapper{cache}
}

func newFastCache() GenericCache {
	cache := fastcache.New(cacheSize)
	return &fastcacheWrapper{cache}
}

func BenchmarkBigCacheGet(b *testing.B) {
	benchmarkGenericGet(b, newBigCache(b))
}

func BenchmarkFastCacheGet(b *testing.B) {
	benchmarkGenericGet(b, newFastCache())
}

func benchmarkGenericGet(b *testing.B, cache GenericCache) {
	rng := rand.New(rand.NewSource(0))

	key := make([]byte, binary.MaxVarintLen64)
	value := make([]byte, cacheEntrySize)

	for i := uint64(0); i < cacheEntryCount; i++ {
		binary.NativeEndian.PutUint64(key, i)
		rng.Read(value[0:])
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

func TestBigCacheCorrectness(t *testing.T) {
	testGenericCorrectness(t, newBigCache(t))
}

func TestFastCacheCorrectness(t *testing.T) {
	testGenericCorrectness(t, newFastCache())
}

func testGenericCorrectness(t *testing.T, cache GenericCache) {
	rng := rand.New(rand.NewSource(0))

	key := make([]byte, binary.MaxVarintLen64)
	value := make([]byte, cacheEntrySize)

	for i := uint64(0); i < cacheEntryCount; i++ {
		binary.NativeEndian.PutUint64(key, i)
		rng.Read(value[0:])
		cache.Set(key, value)
	}

	// Let's reset the random generator
	rng = rand.New(rand.NewSource(0))

	for i := uint64(0); i < cacheEntryCount; i++ {
		binary.NativeEndian.PutUint64(key, i)
		rng.Read(value[0:])
		require.Equal(t, value, cache.Get(key), "failed to verify key %d", i)
	}
}
