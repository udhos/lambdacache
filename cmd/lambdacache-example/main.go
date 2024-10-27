// Package main implements the example.
package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/udhos/lambdacache/lambdacache"
)

func main() {

	debug := os.Getenv("DEBUG") != ""

	if debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	options := lambdacache.Options{
		Debug:    debug,
		Retrieve: getInfo,
		//CleanupInterval: 1, // 1 means always pratically
	}

	cache := lambdacache.New(options)

	for i := range 10 {
		id := i % 2
		key := fmt.Sprintf("key%d", id+1)
		begin := time.Now()
		value, err := cache.Get(key)
		elapsed := time.Since(begin)
		fmt.Printf("key=%s value=%s elap=%v error=%v\n",
			key, value, elapsed, err)
	}
}

func getInfo(key string) (interface{}, time.Duration, error) {
	time.Sleep(100 * time.Millisecond) // adds fake latency
	value := key + ":value"
	const ttl = 2 * time.Second // per-key ttl
	return value, ttl, nil
}
