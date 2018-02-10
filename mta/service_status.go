package mta

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// APIServiceResponse is the root node.
type APIServiceResponse struct {
	ResponseCode int    `xml:"responsecode"`
	Updated      string `xml:"timestamp"`
	Subway       Subway `xml:"subway"`
}

// Subway is the main mode of transport.
type Subway struct {
	Lines []SubwayLine `xml:"line"`
}

// SubwayLine is a line for the subway.
type SubwayLine struct {
	Name   string `xml:"name"`
	Status string `xml:"status"`
	Text   string `xml:"text"`
}

type Status struct {
	Line   string
	Status string
}

type Service struct {
	Updated *time.Time
	Status  []*Status
}

const serviceStatusURL = "http://web.mta.info/status/serviceStatus.txt"

// GetServiceStatus returns the current service status.
func (c *Client) GetServiceStatus() (*Service, error) {
	req, _ := http.NewRequest("GET", serviceStatusURL, nil)
	resp, err := c.httpClient().Do(req)
	defer mustClose(resp.Body)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var service *APIServiceResponse
	err = xml.Unmarshal(body, &service)
	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("America/New_York")
	updated, err := time.ParseInLocation("2/1/2006 3:04:05 PM", service.Updated, loc)
	if err != nil {
		return nil, err
	}
	updated = updated.UTC()

	status := make([]*Status, 0, len(service.Subway.Lines))
	for _, line := range service.Subway.Lines {
		status = append(status, &Status{
			Line:   line.Name,
			Status: strings.Title(line.Status),
		})
	}

	return &Service{Updated: &updated, Status: status}, nil
}
