// Package lambdacache implements a cache for aws lambda go functions.
package lambdacache

import (
	"log/slog"
	"time"
)

// Options define cache settings.
type Options struct {
	// Retrieve is required cache filling function to retrieve a key that is
	// missing from the cache. It must return the tuple "value, TTL, error".
	// Value is value for the key, TTL is how long the key should be kept in
	// cache, and error is used to signal any error that prevented key
	// retrieval.
	Retrieve func(key string) (string, time.Duration, error)

	// If unset, defaults to 5 minutes.
	// Set to negative value (ie -1) to disable clean-up.
	CleanupInterval time.Duration

	// Debug enables debug logging.
	Debug bool

	// Time is a pluggable time source interface for testing.
	Time TimeSource
}

// Cache represents a cache instance.
type Cache struct {
	options     Options
	cache       map[string]entry
	lastCleanup time.Time
}

type entry struct {
	deadline time.Time
	value    string
}

func (e entry) isAlive(now time.Time) bool {
	return e.deadline.After(now)
}

// New creates a cache instance.
func New(options Options) *Cache {

	if options.Retrieve == nil {
		panic("Options.Retrieve is required")
	}

	if options.CleanupInterval == 0 {
		options.CleanupInterval = 5 * time.Minute
	}

	if options.Time == nil {
		options.Time = defaultTime{}
	}

	return &Cache{
		options: options,
		cache:   map[string]entry{},
	}
}

// Get gets value for key from cache.
func (c *Cache) Get(key string) (string, error) {

	now := c.options.Time.Now()

	if c.options.CleanupInterval > 0 && c.options.Time.Since(c.lastCleanup) > c.options.CleanupInterval {
		size := len(c.cache)
		for k, e := range c.cache {
			if !e.isAlive(now) {
				delete(c.cache, k)
			}
		}

		remain := len(c.cache)
		deleted := size - remain
		if c.options.Debug {
			slog.Debug("lambdacache.Cache.Get: cleanup",
				"elapsed", c.options.Time.Since(now),
				"scanned", size,
				"deleted", deleted,
				"remain", remain,
			)
		}

		c.lastCleanup = now
	}

	e, found := c.cache[key]
	if found {
		if e.isAlive(c.options.Time.Now()) {
			return e.value, nil
		}
		delete(c.cache, key)
	}

	v, ttl, errRetrieve := c.options.Retrieve(key)
	if errRetrieve != nil {
		return "", errRetrieve
	}

	e = entry{
		value:    v,
		deadline: c.options.Time.Now().Add(ttl),
	}

	c.cache[key] = e

	return v, nil
}