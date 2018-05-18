package mta

import (
	"testing"
)

func parse(t *testing.T) *parseResult {
	p := Parser{"./testdata/gtfs/stops.txt", "./testdata/gtfs/transfers.txt"}
	res, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func client(t *testing.T) *Client {
	client, err := NewClient(&ClientConfig{
		StopsFilePath:     "./testdata/gtfs/stops.txt",
		TransfersFilePath: "./testdata/gtfs/transfers.txt",
	})
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func TestStationMapping(t *testing.T) {
	var tests = []struct {
		id   string
		name string
	}{
		{"635", "Union Sq - 14 St"},
		{"L03", "Union Sq - 14 St"},
		{"725", "Times Sq - 42 St"},
		{"127", "Times Sq - 42 St"},
		{"A27", "42 St - Port Authority Bus Terminal"},
		{"132", "14 St"},
	}
	res := parse(t)
	for _, tt := range tests {
		id := res.StationMap[tt.id]
		station := res.Stations[id]
		if station.Name != tt.name {
			t.Errorf("station(%v) got %v, want %v", tt.id, station.Name, tt.name)
		}
	}
}

func TestGetClosestStations(t *testing.T) {
	var tests = []struct {
		coordinates *Coordinates
		expected    StationID
	}{
		{&Coordinates{40.7347908, -73.9907299}, "L03"},
		{&Coordinates{40.7376712, -73.992523}, "L03"},
		{&Coordinates{40.7375249, -73.9969781}, "D19"},
		{&Coordinates{40.7387666, -73.9997193}, "132"},
	}

	c := client(t)
	for _, tt := range tests {
		stations := c.GetClosestStations(tt.coordinates, 1)
		if len(stations) != 1 {
			t.Errorf("stations got %v, want %v", len(stations), 1)
		}
		if stations[0].ID != tt.expected {
			t.Errorf("StationID got %v, want %v", stations[0].ID, tt.expected)
		}
	}
}

func TestGetStation(t *testing.T) {
	var tests = []struct {
		id    StationID
		name  string
		error bool
	}{
		{"132", "14 St", false},
		{"L03", "Union Sq - 14 St", false},
		{"foo", "", true},
	}

	c := client(t)
	for _, tt := range tests {
		station, err := c.GetStation(tt.id)
		if err == nil && tt.error {
			t.Error("GetStation expected no error")
		}

		if !tt.error && station.Name != tt.name {
			t.Errorf("station.Name got %v, want %v", station.Name, tt.name)
		}
	}
}

func TestGetStations(t *testing.T) {
	c := client(t)
	expected := 414
	if len(c.GetStations()) != expected {
		t.Errorf("got %v want %v", len(c.GetStations()), expected)
	}
}
