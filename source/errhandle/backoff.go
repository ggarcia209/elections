// Package errhandle contains objects and methods for handling errors
package errhandle

import (
	"math"
	"time"
)

// ExponentialBackoffConfig enables changes to the exponential backoff algorithm parameters
// Attempt, Elapsed, MaxRetiresReached should always be initialized to 0, 0, false
type ExponentialBackoffConfig struct {
	Base              float64
	Cap               float64
	Attempt           float64
	Elapsed           float64
	MaxRetriesReached bool
}

// DefaultFailConfig is the default configuration for the exponential backoff alogrithm
// with a base wait time of 50 miliseconds, and max wait time of 1 minute (60000 ms)
var DefaultFailConfig = &ExponentialBackoffConfig{50, 60000, 0, 0, false}

// ExponentialBackoff implements the exponential backoff algorithm for request retries
// and returns true when the max number of retries has been reached (fc.Elapsed > fc.Cap)
func (fc *ExponentialBackoffConfig) ExponentialBackoff() {
	fc.Attempt += 1.0
	wait := fc.Base * math.Pow(2.0, fc.Attempt)

	if fc.Elapsed+wait > fc.Cap {
		time.Sleep(time.Duration(wait - (wait + fc.Elapsed - fc.Cap)))
		fc.MaxRetriesReached = true // max retries reached
	}

	time.Sleep(time.Duration(wait) * time.Millisecond)
	fc.Elapsed += wait
}

// Reset resets Attempt and Elapsed fields
func (fc *ExponentialBackoffConfig) Reset() {
	fc.Attempt = 0
	fc.Elapsed = 0
	fc.MaxRetriesReached = false
}