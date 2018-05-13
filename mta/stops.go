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
	TripID  string
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

func (s Schedule) contains(u *Update) int {
	for k, update := range s {
		if u.TripID == update.TripID {
			return k
		}
	}
	return -1
}

func cleanupSchedule(s Schedule) Schedule {
	now := time.Now().UTC()
	tripIDs := make(map[string]struct{})
	y := s[:0]
	for _, n := range s {
		if _, ok := tripIDs[n.TripID]; !ok {
			if n.Arrival.After(now) {
				y = append(y, n)
			}
			tripIDs[n.TripID] = struct{}{}
		}
	}
	return y
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
