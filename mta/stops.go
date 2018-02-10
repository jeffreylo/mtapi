package mta

import (
	"time"
)

const defaultStopsFile = "/etc/mta/gtfs/stops.txt"
const defaultTransfersFile = "/etc/mta/gtfs/transfers.txt"

type StopID string

// Stops keys stops by id.
type Stops map[StopID]*Stop

// Coordinates represents a point on the Earth's surface.
type Coordinates struct {
	Lat float64
	Lon float64
}

// Update is a truncation of the GTFS spec.
type Update struct {
	RouteID string
	Arrival *time.Time
	Delay   int32
}

// ScheduleByArrival sorts Schedule by arrival time.
type ScheduleByArrival Schedule

func (s ScheduleByArrival) Len() int {
	return len(s)
}
func (s ScheduleByArrival) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ScheduleByArrival) Less(i, j int) bool {
	return s[i].Arrival.Before(*s[j].Arrival)
}

// Schedule defines the times for a route.
type Schedule []*Update

func (s Schedule) contains(u *Update) bool {
	for _, update := range s {
		inputArrival := u.Arrival.Truncate(60 * time.Second)
		updateArrival := (*update.Arrival).Truncate(60 * time.Second)
		if inputArrival.Sub(updateArrival).Seconds() <= 30 && u.RouteID == update.RouteID {
			return true
		}
	}
	return false
}

func cleanupSchedule(s Schedule) {
	now := time.Now()
	y := s[:0]
	for _, n := range s {
		if n.Arrival.After(now) {
			y = append(y, n)
		}
	}
}

// Direction represents the vector.
type Direction string

// Stop describes a station in the system.
type Stop struct {
	ID          StopID
	Name        string
	Coordinates *Coordinates
	Schedules   map[Direction]Schedule
	Updated     *time.Time
}

// GetStop returns a Stop by ID. If the stop is not known, a placeholder
// stop will be created.
func (c *Client) GetStop(stopID StopID) *Stop {
	stop, ok := c.stops[stopID]
	if !ok {
		return nil
	}
	return stop
}
