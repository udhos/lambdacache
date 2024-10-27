package lambdacache

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestCache(t *testing.T) {

	const ttl = 2 * time.Second // per-key ttl

	var counter int
	var cacheMisses int

	retrieve := func(key string) (interface{}, time.Duration, error) {
		cacheMisses++
		time.Sleep(100 * time.Millisecond) // adds fake latency
		value := fmt.Sprintf("%s.%d", key, counter)
		return value, ttl, nil
	}

	clock := newTime(time.Now())

	options := Options{
		Retrieve:        retrieve,
		Time:            clock,
		CleanupInterval: ttl,
	}

	cache := New(options)

	key1 := "key1"

	//
	// First query for key1
	//

	value1, errGet1 := cache.Get(key1)
	if errGet1 != nil {
		t.Errorf("key1 get error: %v", errGet1)
	}

	if value1 != "key1.0" {
		t.Errorf("key1 value error: %s", value1)
	}

	if cacheMisses != 1 {
		t.Errorf("expecting 1 cache miss, got %d", cacheMisses)
	}

	if len(cache.cache) != 1 {
		t.Errorf("after key1 1st query, cache should have 1 entry, but found %d",
			len(cache.cache))
	}

	//
	// Second query for key1
	//

	// counter is increased to 1, but value from cache should still be key1.0
	counter++

	value1, errGet1 = cache.Get(key1)
	if errGet1 != nil {
		t.Errorf("key1 get error: %v", errGet1)
	}

	if value1 != "key1.0" {
		t.Errorf("key1 value error: %s", value1)
	}

	if cacheMisses != 1 {
		t.Errorf("expecting 1 cache miss, got %d", cacheMisses)
	}

	if len(cache.cache) != 1 {
		t.Errorf("after key1 2nd query, cache should have 1 entry, but found %d",
			len(cache.cache))
	}

	//
	// Third query for key1
	//

	// counter is 1, and value from cache should now be key1.1

	// advance time to force cache miss
	clock.SetNow(clock.now.Add(ttl))

	value1, errGet1 = cache.Get(key1)
	if errGet1 != nil {
		t.Errorf("key1 get error: %v", errGet1)
	}

	if value1 != "key1.1" {
		t.Errorf("key1 value error: %s", value1)
	}

	if cacheMisses != 2 {
		t.Errorf("expecting 2 cache misses, got %d", cacheMisses)
	}

	if len(cache.cache) != 1 {
		t.Errorf("after key1 3rd query, cache should have 1 entry, but found %d",
			len(cache.cache))
	}

	key2 := "key2"

	//
	// First query for key2
	//

	value2, errGet2 := cache.Get(key2)
	if errGet2 != nil {
		t.Errorf("key2 get error: %v", errGet2)
	}

	if value2 != "key2.1" {
		t.Errorf("key2 value error: %s", value2)
	}

	if cacheMisses != 3 {
		t.Errorf("expecting 3 cache misses, got %d", cacheMisses)
	}

	if len(cache.cache) != 2 {
		t.Errorf("after key2 1st query, cache should have 2 entries, but found %d",
			len(cache.cache))
	}
}

func TestCacheError(t *testing.T) {

	retrieve := func(_ string) (interface{}, time.Duration, error) {
		return "", 0, errors.New("retrieve error")
	}

	options := Options{
		Retrieve: retrieve,
	}

	cache := New(options)

	key1 := "key1"

	_, errGet1 := cache.Get(key1)
	if errGet1 == nil {
		t.Errorf("expecting retrieve error but got success")
	}
}

func TestCacheInt(t *testing.T) {

	const ttl = 2 * time.Second // per-key ttl

	retrieve := func(key string) (interface{}, time.Duration, error) {
		value, err := strconv.Atoi(key)
		return 10 * value, ttl, err
	}

	options := Options{
		Retrieve: retrieve,
	}

	cache := New(options)

	key1 := "2"

	value, errGet1 := cache.Get(key1)
	if errGet1 != nil {
		t.Errorf("key1 get error: %v", errGet1)
	}

	if value != 20 {
		t.Errorf("expecting key1 value=20, but got %d", value)
	}
}

func TestCacheIntError(t *testing.T) {

	const ttl = 2 * time.Second // per-key ttl

	retrieve := func(key string) (interface{}, time.Duration, error) {
		value, err := strconv.Atoi(key)
		return 10 * value, ttl, err
	}

	options := Options{
		Retrieve: retrieve,
	}

	cache := New(options)

	key1 := "2"

	value, errGet1 := cache.Get(key1)
	if errGet1 != nil {
		t.Errorf("key1 get error: %v", errGet1)
	}

	if value == "20" {
		t.Errorf("expecting int, but got string")
	}
}

type customTime struct {
	now time.Time
}

func newTime(now time.Time) *customTime {
	return &customTime{now: now}
}

func (t *customTime) SetNow(now time.Time) {
	t.now = now
}

// Now returns the current local time.
func (t customTime) Now() time.Time {
	return t.now
}

// Since returns the time elapsed since u.
func (t customTime) Since(u time.Time) time.Duration {
	return t.now.Sub(u)
}
