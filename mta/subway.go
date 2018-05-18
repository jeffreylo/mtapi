package mta

import (
	"time"

	"github.com/kyroy/kdtree/points"
	"github.com/pkg/errors"
)

const defaultStopsFile = "/etc/mta/gtfs/stops.txt"
const defaultTransfersFile = "/etc/mta/gtfs/transfers.txt"

// Direction represents the vector.
type Direction string

// StationID is the identifier for a station.
type StationID string

// Stations is a map of StationIDs to station.
type Stations map[StationID]*Station

// Station is a point in the transit system.
type Station struct {
	ID          StationID
	Name        string
	Coordinates *Coordinates
	Arrivals    map[Direction][]*Arrival
	Updated     *time.Time
}

// Arrival is a truncation of the GTFS spec.
type Arrival struct {
	TripID  string
	RouteID string
	Time    *time.Time
}

// Coordinates represents a point on the Earth's surface.
type Coordinates struct {
	Lat float64
	Lon float64
}

// ByArrivalTime sorts Arrivals by arrival time.
type ByArrivalTime []*Arrival

func (s ByArrivalTime) Len() int {
	return len(s)
}
func (s ByArrivalTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByArrivalTime) Less(i, j int) bool {
	var t time.Time
	if s[j].Time != nil {
		t = *s[j].Time
	}
	return s[i].Time.Before(t)
}

func cleanupArrivals(s []*Arrival) []*Arrival {
	now := time.Now().UTC()
	tripIDs := make(map[string]struct{})
	y := s[:0]
	for _, n := range s {
		if _, ok := tripIDs[n.TripID]; !ok {
			if n.Time.After(now) {
				y = append(y, n)
			}
			tripIDs[n.TripID] = struct{}{}
		}
	}
	return y
}

const maxStations = 5

// GetStations returns all stations.
func (c *Client) GetStations() Stations { return c.stations }

var errStationNotFound = errors.New("station not found")

// GetStationByStopID returns the station, i.e., an aggregation of GTFS
// stops, for the GTFS stop id.
func (c *Client) GetStationByStopID(id string) (*Station, error) {
	stationID, ok := c.stops[id]
	if !ok {
		return nil, errStationNotFound
	}
	station, ok := c.stations[stationID]
	if !ok {
		return nil, errStationNotFound
	}
	return station, nil
}

// GetStation returns a station
func (c *Client) GetStation(id StationID) (*Station, error) {
	s, ok := c.stations[id]
	if !ok {
		return nil, errStationNotFound
	}
	return s, nil
}

// GetClosestStations returns the closest stations for the given coordinates.
func (c *Client) GetClosestStations(v *Coordinates, numStations int) []*Station {
	if numStations >= maxStations {
		numStations = maxStations
	} else if numStations <= 0 {
		numStations = 1
	}
	results := c.tree.KNN(&points.Point{Coordinates: []float64{v.Lat, v.Lon}}, numStations)
	stations := make([]*Station, 0, len(results))
	for _, v := range results {
		point := v.(*points.Point)
		station, err := c.GetStation(point.Data.(StationID))
		if err != nil {
			continue
		}
		stations = append(stations, station)
	}
	return stations
}
