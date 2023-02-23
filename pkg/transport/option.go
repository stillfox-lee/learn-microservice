package transport

import "time"

// create an options for Ratelimit, Circuitbreaker, Retry and Timeout.
type TrafficControlOptions struct {
	QPS            int           // rate limit per second
	Burst          int           // request concurrency number
	BreakerTimeout time.Duration // circuitbreaker timeout
	RetryMax       int           // retry max times
	RetryTimeout   time.Duration // retry timeout
}
