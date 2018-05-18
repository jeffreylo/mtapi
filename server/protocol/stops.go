package protocol

import (
	"strings"
	"time"

	"github.com/jeffreylo/mtapi/mta"
)

type Coordinates struct {
	Lat float64
	Lon float64
}

type Arrival struct {
	TripID  string
	Time    *time.Time
	RouteID string
}

type Arrivals map[mta.Direction][]*Arrival

type Station struct {
	ID          string
	Name        string
	Coordinates *Coordinates
	Arrivals    map[mta.Direction][]*Arrival `json:",omitempty"`
	Updated     *time.Time                   `json:",omitempty"`
}

func (p *Protocol) Arrivals(v map[mta.Direction][]*mta.Arrival) Arrivals {
	w := make(Arrivals)
	for d, s := range v {
		vv := make([]*Arrival, 0, len(s))
		for _, u := range s {
			routeID := u.RouteID
			if strings.HasSuffix(routeID, "S") {
				routeID = "S"
			}
			vv = append(vv, &Arrival{
				TripID:  u.TripID,
				Time:    u.Time,
				RouteID: routeID,
			})
		}
		w[d] = vv
	}
	return w
}

func (p *Protocol) Station(v *mta.Station) *Station {
	return &Station{
		ID:   string(v.ID),
		Name: v.Name,
		Coordinates: &Coordinates{
			Lat: v.Coordinates.Lat,
			Lon: v.Coordinates.Lon,
		},
		Arrivals: p.Arrivals(v.Arrivals),
	}
}

func (p *Protocol) Stations(stations mta.Stations) []*Station {
	result := make([]*Station, 0, len(stations))
	for _, station := range stations {
		result = append(result, &Station{
			ID:   string(station.ID),
			Name: station.Name,
			Coordinates: &Coordinates{
				Lat: station.Coordinates.Lat,
				Lon: station.Coordinates.Lon,
			},
		})
	}
	return result
}
