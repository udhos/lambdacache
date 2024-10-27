package lambdacache

import "time"

// TimeSource is a pluggable time source interface for testing.
type TimeSource interface {
	Now() time.Time
	Since(t time.Time) time.Duration
}

// defaultTime provides default time source implementation.
type defaultTime struct{}

// Now returns the current local time.
func (t defaultTime) Now() time.Time {
	return time.Now()
}

// Since returns the time elapsed since u.
func (t defaultTime) Since(u time.Time) time.Duration {
	return t.Now().Sub(u)
}
