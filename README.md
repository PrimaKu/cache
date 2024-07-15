# Cache Manager

A Go library for managing cache with Redis.

## Installation

```bash
go get github.com/PrimaKu/cache
```

## Setup

```go
cacheManager, err := cache.NewCacheManager(redisOptions) // *redis.Options
if err != nil {
  log.Fatalf("Failed to start caching with: %v", err)
}
```

## Get Cache
```go
data := cacheManager.Get(context, "DATA_KEY")
```

## Set Cache
```go
err := cacheManager.Set(context, "DATA_KEY", "value to cache")
```

## Delete Cache
```go
err := cacheManager.Del(context, "DATA_KEY")
```