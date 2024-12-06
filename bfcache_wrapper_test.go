package bfcache_test

import (
	"unsafe"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/allegro/bigcache/v3"
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

func bytesToString(b []byte) string {
	data := unsafe.SliceData(b)
	return unsafe.String(data, len(b))
}
