package mta

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

// Service is the root node.
type Service struct {
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
	Date   string `xml:"date"`
	Time   string `xml:"time"`
}

const serviceStatusURL = "http://web.mta.info/status/serviceStatus.txt"

// GetServiceStatus returns the current service status.
func GetServiceStatus() (*Service, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", serviceStatusURL, nil)
	resp, err := client.Do(req)
	defer mustClose(resp.Body)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var service *Service
	err = xml.Unmarshal(body, &service)
	if err != nil {
		return nil, err
	}
	return service, nil
}
