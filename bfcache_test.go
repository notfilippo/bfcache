package bfcache_test

import (
	"testing"

	"github.com/notfilippo/bfcache"
	"github.com/stretchr/testify/require"
)

func TestSetGet(t *testing.T) {
	cache := bfcache.New()
	cache.Set([]byte("hello"), []byte("world"))

	value := cache.Get([]byte("hello"))
	require.Equal(t, []byte("world"), value)
}
