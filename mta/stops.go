package mta

import (
	"os"
	"regexp"
	"time"

	"github.com/gocarina/gocsv"
)

const defaultStopsFile = "/etc/mta/gtfs/stops.txt"

// Stops keys stops by id.
type Stops map[string]*Stop

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
	ID          string
	Name        string
	Coordinates *Coordinates
	Schedules   map[Direction]Schedule
	Updated     *time.Time
}

// GetStop returns a Stop by ID. If the stop is not known, a placeholder
// stop will be created.
func (c *Client) GetStop(stopID string) *Stop {
	stop, ok := c.stops[stopID]
	if !ok {
		return nil
	}
	return stop
}

const stopRegex = "(?P<ID>.*)(?P<Direction>[NS])"

// ParseStopsFile reads and loads the stops data.
func ParseStopsFile(path string) (Stops, error) {
	re := regexp.MustCompile(stopRegex)
	type result struct {
		ID   string  `csv:"stop_id"`
		Name string  `csv:"stop_name"`
		Lat  float64 `csv:"stop_lat"`
		Lon  float64 `csv:"stop_lon"`
		Type int     `csv:"location_type"`
	}

	if path == "" {
		path = defaultStopsFile
	}

	s, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	var stops []*result
	if err := gocsv.UnmarshalFile(s, &stops); err != nil {
		panic(err)
	}
	v := make(Stops, len(stops))
	for _, stop := range stops {
		if stop.Type == 0 {
			m := re.FindAllStringSubmatch(stop.ID, -1)[0]
			v[m[1]] = &Stop{
				ID:   m[1],
				Name: stop.Name,
				Coordinates: &Coordinates{
					Lat: stop.Lat,
					Lon: stop.Lon,
				},
				Schedules: make(map[Direction]Schedule),
			}
		}
	}
	return v, nil
}
