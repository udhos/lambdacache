// Package main implements the example.
package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/udhos/lambdacache/lambdacache"
)

// create cache in lambda function GLOBAL context
var cache = newCache()

func main() {

	for i := range 10 {
		id := i % 2
		key := fmt.Sprintf("key%d", id+1)
		begin := time.Now()

		// query cache like this in lambda function HANDLER context
		value, err := cache.Get(key)

		elapsed := time.Since(begin)
		fmt.Printf("key=%s value=%s elap=%v error=%v\n",
			key, value, elapsed, err)
	}
}

func newCache() *lambdacache.Cache {
	debug := os.Getenv("DEBUG") != ""

	if debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	options := lambdacache.Options{
		Debug:    debug,
		Retrieve: getInfo,
		//CleanupInterval: 1, // 1 means always pratically
	}

	return lambdacache.New(options)
}

func getInfo(key string) (interface{}, time.Duration, error) {
	time.Sleep(100 * time.Millisecond) // adds fake latency
	value := key + ":value"
	const ttl = 2 * time.Second // per-key ttl
	return value, ttl, nil
}
