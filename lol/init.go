package lol

import "sync"

var (
	defaultClient *Client
	one           sync.Once
)

const (
	NA1 = "NA1"
)

// DefaultClient returns the default client
func DefaultClient() *Client {
	one.Do(func() {
		c, err := NewClient()
		if err != nil {
			panic(err)
		}
		defaultClient = c
	})
	return defaultClient
}
