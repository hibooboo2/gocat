package lol

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

type Client struct {
	client            *http.Client
	cache             lolStorer
	requestsMade      *int64
	requestsSucceeded *int64
	Debug             bool
}

func NewClient() (*Client, error) {
	cache, err := NewLolMongo()
	if err != nil {
		return nil, err
	}
	var x, y int64
	return &Client{
		client: &http.Client{
			Timeout: time.Second * 5,
		},
		requestsMade:      &x,
		requestsSucceeded: &y,
		cache:             cache,
	}, nil
}

func (c *Client) GetCache() lolStorer {
	return c.cache
}

func (c *Client) Get(url string) (*http.Response, error) {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// if os.Getenv("X_Riot_Token") != "" {
	// 	r.Header.Add("X-Riot-Token", os.Getenv("X_Riot_Token"))
	// }
	resp, err := c.client.Do(r)
	if c.Debug {
		fmt.Fprintf(os.Stdout, "\t\t\t\t\t\t\t\t\t\t\tRequests Made: %d Requests Succeeded: %d\r", atomic.AddInt64(c.requestsMade, 1), atomic.LoadInt64(c.requestsSucceeded))
	}
	if err != nil {
		return resp, err
	}
	switch resp.StatusCode {
	case http.StatusTooManyRequests:
		time.Sleep(time.Second * 2)
		if c.Debug {
			fmt.Fprintf(os.Stdout, "debug: slow down charlie.\r")
		}
		return c.Get(url)
	case http.StatusNotFound:
		if c.Debug {
			log.Println("\nerr: not found", url)
		}
		return resp, err
	}
	atomic.AddInt64(c.requestsSucceeded, 1)
	return resp, err
}

func (c *Client) Close() {
	c.cache.Close()
}
