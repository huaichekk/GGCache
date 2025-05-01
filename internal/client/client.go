package client

import (
	"GGCache/internal/cache/eviction"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Client struct {
	addr string
}

func NewClient(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

func (c *Client) Get(group, key string) (eviction.ByteView, bool) {
	url := fmt.Sprintf("http://%s/%s/%s", c.addr, group, key)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return eviction.ByteView{}, false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println(resp.StatusCode)
		return eviction.ByteView{}, false
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return eviction.ByteView{}, false
	}
	return eviction.ByteView{B: body}, true
}
