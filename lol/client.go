package lol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	BR   RegionEndPoint = "br1.api.riotgames.com"
	EUNE RegionEndPoint = "eun1.api.riotgames.com"
	EUW  RegionEndPoint = "euw1.api.riotgames.com"
	JP   RegionEndPoint = "jp1.api.riotgames.com"
	KR   RegionEndPoint = "kr.api.riotgames.com"
	LAN  RegionEndPoint = "la1.api.riotgames.com"
	LAS  RegionEndPoint = "la2.api.riotgames.com"
	NA   RegionEndPoint = "na1.api.riotgames.com"
	OCE  RegionEndPoint = "oc1.api.riotgames.com"
	TR   RegionEndPoint = "tr1.api.riotgames.com"
	RU   RegionEndPoint = "ru.api.riotgames.com"
	PBE  RegionEndPoint = "pbe1.api.riotgames.com"
)

type RegionEndPoint string

type client struct {
	client            *http.Client
	baseURL           RegionEndPoint
	requests          map[string]string
	cache             lolStorer
	requestsMade      *int64
	requestsSucceeded *int64
	requestLock       sync.RWMutex
}

// Debug weather or not debug is enabled for the riot package.
var Debug bool

// NewClient returns a new client for performing operations using riots api.
func NewClient(region RegionEndPoint) (RiotClient, error) {
	// cache, err := NewLolMongo("dev.jhrb.us", 27217)
	cache, err := NewLolMongo("", 0)
	// cache, err := NewLolMongo("192.168.1.170", 27027)
	if err != nil {
		return nil, err
	}
	var x, y int64
	return &client{
		client: &http.Client{
			Timeout: time.Second * 5,
		},
		requestsMade:      &x,
		requestsSucceeded: &y,
		cache:             cache,
		baseURL:           region,
		requests:          make(map[string]string),
	}, nil
}

func (c *client) GetCache() lolStorer {
	return c.cache
}

func (c *client) GetObjRiot(url string, val interface{}) error {
	url = path.Join(string(c.baseURL), url)
	body, err := c.GetBody(url)
	if err != nil {
		return err
	}
	buff := &bytes.Buffer{}
	io.Copy(buff, body)
	cp := buff.String()
	err = json.NewDecoder(buff).Decode(val)
	if err != nil {
		logger.Println("trace: body: ", cp)
	}
	logger.Println("trace: body: ", cp)

	return err
}

func (c *client) GetBody(url string) (io.Reader, error) {
	c.requestLock.RLock()
	body, ok := c.requests[url]
	c.requestLock.RUnlock()
	if ok {
		return strings.NewReader(body), nil
	}
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buff := &bytes.Buffer{}
	io.Copy(buff, resp.Body)
	c.requestLock.Lock()
	c.requests[url] = buff.String()
	c.requestLock.Unlock()
	return buff, nil
}

func (c *client) Get(url string) (*http.Response, error) {
	c.requestLock.Lock()
	defer c.requestLock.Unlock()
	if !strings.HasPrefix(url, "http") {
		url = fmt.Sprintf("https://%s", url)
	}
	logger.Printf("trace: GET: %s", url)
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if os.Getenv("X_Riot_Token") != "" {
		r.Header.Add("X-Riot-Token", os.Getenv("X_Riot_Token"))
	}
	resp, err := c.client.Do(r)
	if Debug {
		fmt.Fprintf(os.Stdout, "\t\t\t\t\t\t\t\t\t\t\tRequests Made: %d Requests Succeeded: %d\r", atomic.AddInt64(c.requestsMade, 1), atomic.LoadInt64(c.requestsSucceeded))
	}
	if err != nil {
		return resp, err
	}
	switch resp.StatusCode {
	case http.StatusTooManyRequests:
		logger.Println("debug: Headers on 429 request:", resp.Header)
		time.Sleep(time.Second * 2)
		logger.Println("trace: slow down charlie.\r")
		c.requestLock.Unlock()
		resp, err = c.Get(url)
		c.requestLock.Lock()
		return resp, err
	case http.StatusNotFound:
		logger.Println("err: not found", url)
		return nil, fmt.Errorf("err: object not found: %s", url)
	}
	atomic.AddInt64(c.requestsSucceeded, 1)
	return resp, err
}

func (c *client) Close() {
	c.cache.Close()
}
