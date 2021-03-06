package mta

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"time"

	raven "github.com/getsentry/raven-go"
	"github.com/kyroy/kdtree"
)

const (
	feedBaseURL     = "http://datamine.mta.info/mta_esi.php"
	refreshInterval = time.Second * 5
)

// datamine.mta.info/list-of-feeds
var feedIDs = []int{1, 2, 16, 21, 26, 31, 36, 51}

// Client consumes the MTA API.
type Client struct {
	apiKey    string
	client    *http.Client
	ignoreSSL bool
	port      int

	stops    map[string]StationID
	stations Stations
	tree     *kdtree.KDTree
	mtx      *sync.Mutex

	err     chan error
	done    chan struct{}
	updated *time.Time
}

// ClientConfig defines the settings for the MTA client.
type ClientConfig struct {
	APIKey            string
	IgnoreSSL         bool
	Port              int
	StopsFilePath     string
	TransfersFilePath string
}

// NewClient returns a new instance of the MTA client.
func NewClient(cfg *ClientConfig) (*Client, error) {
	parser := &Parser{cfg.StopsFilePath, cfg.TransfersFilePath}
	result, err := parser.Parse()
	if err != nil {
		return nil, err
	}
	c := &Client{
		apiKey:   cfg.APIKey,
		done:     make(chan struct{}),
		err:      make(chan error),
		mtx:      &sync.Mutex{},
		port:     cfg.Port,
		stations: result.Stations,
		stops:    result.StationMap,
		tree:     result.Tree,
	}
	return c, nil
}

// Close terminates the session.
func (c *Client) Close() error {
	c.done <- struct{}{}
	return nil
}

// Work starts the feed fetcher.
func (c *Client) Work() {
	raven.CapturePanic(func() {
		c.refreshFeeds()
		ticker := time.NewTicker(refreshInterval)
		for {
			select {
			case <-ticker.C:
				c.refreshFeeds()
			case <-c.done:
				ticker.Stop()
				return
			}
		}
	}, nil)
}

func (c *Client) refreshFeeds() {
	var wg sync.WaitGroup
	wg.Add(len(feedIDs))
	for _, feedID := range feedIDs {
		go func(feedID int) {
			defer wg.Done()
			c.refreshFeed(feedID)
		}(feedID)
	}
	wg.Wait()
}

func (c *Client) httpClient() *http.Client {
	if c.client == nil {
		if c.ignoreSSL {
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			c.client = &http.Client{Transport: tr}
		} else {
			c.client = http.DefaultClient
		}
	}
	return c.client
}

func (c *Client) getFeedURL(feedID int) string {
	return fmt.Sprintf("%s?&key=%s&feed_id=%d", feedBaseURL, c.apiKey, feedID)
}
