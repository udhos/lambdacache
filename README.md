[![license](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/udhos/lambdacache/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/udhos/lambdacache)](https://goreportcard.com/report/github.com/udhos/lambdacache)
[![Go Reference](https://pkg.go.dev/badge/github.com/udhos/lambdacache.svg)](https://pkg.go.dev/github.com/udhos/lambdacache)

# lambdacache

[lambdacache](https://github.com/udhos/lambdacache) is a Go package for caching data in your AWS Lambda function between sequential invocations. It stores data in the memory of the function [execution context](https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtime-environment.html) that usually spans multiple consecutive invocations.

# Benefits

- Lower load on backend services.
- Lower lambda execution time (thus lower lambda costs).
- Independent of any external cache service (e.g. redis).
- Higher robustness.
- Higher scalability.

# Usage

```golang
// 1. import the package
import "github.com/udhos/lambdacache/lambdacache"

// 2. in lambda function GLOBAL context: create cache
var cache = newCache()

func newCache() *lambdacache.Cache {
    options := lambdacache.Options{
        Debug:    true,
        Retrieve: getInfo,
    }
    return lambdacache.New(options)
}

// 3. in lambda function HANDLER context: query cache
func HandleRequest(ctx context.Context) error {
    // ...
    value, errGet := cache.Get(key)
    // ...
    return nil
}

// getInfo retrieves key value when there is a cache miss
func getInfo(key string) (interface{}, time.Duration, error) {
    const ttl = 5 * time.Minute // per-key TTL
    return "put-retrieved-value-here", ttl, nil
}
```

# Example

See:

[./cmd/lambdacache-example/main.go](./cmd/lambdacache-example/main.go)
