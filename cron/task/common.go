package task

import "time"

// ExponentialBackoffSleep 指数退避休眠
func ExponentialBackoffSleep(current, max time.Duration) time.Duration {
	time.Sleep(current)
	next := current * 2
	if next > max {
		next = max
	}
	return next
}
