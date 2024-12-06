package bfcache_test

import (
	"context"
	"crypto/rand"
	"encoding/binary"
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

func newValue(size uint) []byte {
	buf := make([]byte, size)
	rand.Read(buf)
	return buf
}

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
	keys := make([][]byte, cacheEntryCount)

	for i := uint64(0); i < cacheEntryCount; i++ {
		key := make([]byte, binary.MaxVarintLen64)
		binary.NativeEndian.PutUint64(key, i)
		cache.Set(key, newValue(cacheEntrySize))
		keys[i] = key
	}

	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		var i uint64
		for p.Next() {
			cache.Get(keys[i])
			i += 1
		}
	})
}
