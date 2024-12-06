package showdown_test

import (
	"sync"
	"unsafe"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/allegro/bigcache/v3"
	"github.com/notfilippo/bfcache"
)

type GenericCache interface {
	Set(key []byte, value []byte)
	Get(key []byte) []byte
}

var (
	_ GenericCache = (*bigcacheWrapper)(nil)
	_ GenericCache = (*bigcacheWrapper)(nil)
)

type bigcacheWrapper struct {
	inner *bigcache.BigCache
}

func (b *bigcacheWrapper) Set(key []byte, value []byte) {
	b.inner.Set(bytesToString(key), value)
}

func (b *bigcacheWrapper) Get(key []byte) []byte {
	v, _ := b.inner.Get(bytesToString(key))
	return v
}

type fastcacheWrapper struct {
	inner *fastcache.Cache
}

func (b *fastcacheWrapper) Set(key []byte, value []byte) {
	b.inner.SetBig(key, value)
}

func (b *fastcacheWrapper) Get(key []byte) []byte {
	return b.inner.GetBig(nil, key)
}

type syncmapWrapper struct {
	inner sync.Map
}

func (b *syncmapWrapper) Set(key []byte, value []byte) {
	b.inner.Store(bytesToString(key), value)
}

func (b *syncmapWrapper) Get(key []byte) []byte {
	value, ok := b.inner.Load(bytesToString(key))
	if !ok {
		return nil
	}
	return value.([]byte)
}

type bfcacheWrapper struct {
	inner *bfcache.Cache
}

func (b *bfcacheWrapper) Set(key []byte, value []byte) {
	b.inner.Set(key, value)
}

func (b *bfcacheWrapper) Get(key []byte) []byte {
	return b.inner.Get(key)
}

func bytesToString(b []byte) string {
	data := unsafe.SliceData(b)
	return unsafe.String(data, len(b))
}
