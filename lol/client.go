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
	client       *http.Client
	cache        *lolCache
	requestsMade *int64
}

func NewClient() *Client {
	x := int64(0)
	return &Client{
		client: &http.Client{
			Timeout: time.Second * 5,
		},
		requestsMade: &x,
	}
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
	fmt.Fprintf(os.Stdout, "Requests Made: %d\r", atomic.AddInt64(c.requestsMade, 1))
	if err != nil {
		return resp, err
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		time.Sleep(time.Second * 2)
		log.Println("debug: slow down charlie.")
		return c.Get(url)
	}
	return resp, err
}
