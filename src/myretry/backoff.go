package myretry

import (
	"time"

	"github.com/sethvargo/go-retry"
)

func NewBackofFactory() func() retry.Backoff {
	return func() retry.Backoff {
		r := retry.NewExponential(1 * time.Second)
		r = retry.WithJitter(5*time.Second, r)
		r = retry.WithMaxRetries(5, r)

		return r
	}
}
