package lol

import "sync"

var (
	defaultClient RiotClient
	one           sync.Once
)

const (
	NA1 = "NA1"
)

// DefaultClient returns the default client
func DefaultClient() RiotClient {
	one.Do(func() {
		c, err := NewClient(NA)
		if err != nil {
			panic(err)
		}
		defaultClient = c
	})
	return defaultClient
}
