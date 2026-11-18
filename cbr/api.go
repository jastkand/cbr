package cbr

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	dateFormat = "02/01/2006"
	cacheTTL   = time.Hour
)

const fakeUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"

type Client struct {
	httpClient *http.Client

	mu             sync.Mutex
	cacheKey       string
	cacheResp      *DailyResponse
	cacheFetchedAt time.Time
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{httpClient: httpClient}
}

func (c *Client) Daily(date time.Time) (*DailyResponse, error) {
	key := date.Format(dateFormat)

	if resp, ok := c.cacheGet(key, cacheTTL); ok {
		return resp, nil
	}

	resp, err := c.fetch(date)
	if err != nil {
		// Rates for a past date don't change, so stale data beats an error.
		if resp, ok := c.cacheGet(key, 0); ok {
			return resp, nil
		}
		return nil, err
	}

	c.cachePut(key, resp)
	return resp, nil
}

func (c *Client) fetch(date time.Time) (*DailyResponse, error) {
	url := fmt.Sprintf("https://www.cbr.ru/scripts/XML_daily.asp?date_req=%s", date.Format(dateFormat))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", fakeUserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var daily DailyResponse
	err = xmlDecoder(body).Decode(&daily)
	if err != nil {
		return nil, err
	}
	return &daily, nil
}

func (c *Client) cacheGet(key string, ttl time.Duration) (*DailyResponse, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cacheKey == key && time.Since(c.cacheFetchedAt) <= ttl {
		return c.cacheResp, true
	}
	return nil, false
}

func (c *Client) cachePut(key string, resp *DailyResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cacheKey, c.cacheResp, c.cacheFetchedAt = key, resp, time.Now()
}
