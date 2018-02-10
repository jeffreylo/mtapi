package protocol

import (
	"time"

	"github.com/jeffreylo/mtapi/mta"
)

type Coordinates struct {
	Lat float64
	Lon float64
}

type Update struct {
	Arrival *time.Time
	RouteID string
}

type Schedules map[mta.Direction][]*Update

type Station struct {
	ID          string
	Name        string
	Coordinates *Coordinates
	StopIDs     []mta.StopID `json:",omitempty"`

	Schedules Schedules  `json:",omitempty"`
	Updated   *time.Time `json:",omitempty"`
}

func (p *Protocol) Schedules(v map[mta.Direction]mta.Schedule) Schedules {
	w := make(Schedules)
	for d, s := range v {
		vv := make([]*Update, 0, len(s))
		for _, u := range s {
			vv = append(vv, &Update{
				Arrival: u.Arrival,
				RouteID: u.RouteID,
			})
		}
		w[d] = vv
	}
	return w
}

func (p *Protocol) Station(v *mta.Station, schedules map[mta.Direction]mta.Schedule, updated *time.Time) *Station {
	return &Station{
		ID:   string(v.ID),
		Name: v.Name,
		Coordinates: &Coordinates{
			Lat: v.Coordinates.Lat,
			Lon: v.Coordinates.Lon,
		},
		Schedules: p.Schedules(schedules),
		Updated:   updated,
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
			StopIDs: station.StopIDs(),
		})
	}
	return result
}
