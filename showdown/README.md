# Showdown 🤺

This folder contains benchmarks for 3 cache libraries:

- [`bfcache`](https://github.com/notfilippo/bfcache) (this repo)
- [`fastcache`](https://github.com/VictoriaMetrics/fastcache)
- [`bigcache`](https://github.com/allegro/bigcache)

The benchmark is representative of the particular environment the author was
trying to optimize against:

- Big cache size (~ 10GBs)
- High number of entries (~ 9M)
- Mid entry size (avg 800B, but could go over 64KB)
