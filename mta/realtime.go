package mta

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/gtfs-realtime-bindings/golang/gtfs"
	"github.com/pkg/errors"
)

const stopRegex = "(?P<ID>.*)(?P<Direction>[NS])"

func (c *Client) refreshFeed(feedID int) {
	re := regexp.MustCompile(stopRegex)
	req, _ := http.NewRequest("GET", c.getFeedURL(feedID), nil)
	resp, err := c.httpClient().Do(req)
	if err != nil {
		log.Print(errors.Wrap(err, "mta: request failed"))
		return
	}
	defer mustClose(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(errors.Wrap(err, "mta: read failed"))
		return
	}

	feed := gtfs.FeedMessage{}
	err = proto.Unmarshal(body, &feed)
	if err != nil {
		msg := err.Error()
		if !strings.HasPrefix(msg, "proto") && !strings.HasPrefix(msg, "unexpected EOF") && !strings.HasPrefix(msg, "bad wiretype") {
			log.Print(errors.Wrap(err, "mta: unmarshal failed"))
		}
		return
	}

	now := time.Now().UTC()
	for _, entity := range feed.Entity {
		tripUpdate := entity.GetTripUpdate()
		if tripUpdate == nil {
			continue
		}

		trip := tripUpdate.GetTrip()
		stopTimeUpdates := tripUpdate.GetStopTimeUpdate()
		for _, update := range stopTimeUpdates {
			stopID := update.GetStopId()
			m := re.FindStringSubmatch(stopID)
			if m[1] == "" {
				continue
			}
			stop := c.GetStop(StopID(m[1]))
			if stop == nil {
				continue
			}

			direction := Direction(m[2])
			arrival := update.GetArrival()
			if arrival != nil {
				arrivalTime := time.Unix(arrival.GetTime(), 0).UTC()
				update := &Update{
					RouteID: trip.GetRouteId(),
					Delay:   arrival.GetDelay(),
					Arrival: &arrivalTime,
				}
				if arrivalTime.After(now) && !stop.Schedules[direction].contains(update) {
					stop.Schedules[direction] = append(stop.Schedules[direction], update)
					cleanupSchedule(stop.Schedules[direction])
					sort.Sort(ScheduleByArrival(stop.Schedules[direction]))
					stop.Updated = &now
				}
			}
		}
	}
}

func mustClose(closer io.ReadCloser) {
	if err := closer.Close(); err != nil {
		log.Panic(err)
	}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
