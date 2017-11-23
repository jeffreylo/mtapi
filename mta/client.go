package mta

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	feedBaseURL     = "http://datamine.mta.info/mta_esi.php"
	refreshInterval = time.Second * 5
)

// datamine.mta.info/list-of-feeds
var feedIDs = []int{1, 2, 16, 21, 26, 31}

// Client consumes the MTA API.
type Client struct {
	apiKey    string
	client    *http.Client
	ignoreSSL bool
	port      int

	stops   Stops
	err     chan error
	done    chan struct{}
	updated *time.Time
}

// ClientConfig defines the settings for the MTA client.
type ClientConfig struct {
	APIKey        string
	IgnoreSSL     bool
	Port          int
	StopsFilePath string
}

// NewClient returns a new instance of the MTA client.
func NewClient(cfg *ClientConfig) (*Client, error) {
	stops, err := ParseStopsFile(cfg.StopsFilePath)
	if err != nil {
		return nil, err
	}
	c := &Client{
		apiKey: cfg.APIKey,
		stops:  stops,
		port:   cfg.Port,
		err:    make(chan error),
		done:   make(chan struct{}),
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

func (c *Client) stopHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stopID := req.Form.Get("id")
	if stopID == "" {
		stops := make(map[string]interface{})
		for _, stop := range c.stops {
			stops[stop.ID] = map[string]interface{}{
				"Name":        stop.Name,
				"Coordinates": stop.Coordinates,
			}
		}
		v, err := json.Marshal(stops)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(v)
		return
	}

	stop := c.GetStop(stopID)
	if stop == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	v, err := json.Marshal(map[string]interface{}{
		"Stop": stop,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(v)
}

func (c *Client) statusHandler(w http.ResponseWriter, req *http.Request) {
	s, err := GetServiceStatus()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	v, err := json.Marshal(map[string]interface{}{
		"Service": s,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(v)
}

// Serve returns an http server.
func (c *Client) Serve() error {
	http.HandleFunc("/stops", c.stopHandler)
	http.HandleFunc("/status", c.statusHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", c.port), nil)
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

// GetStops returns all stops.
func (c *Client) GetStops() Stops {
	return c.stops
}
