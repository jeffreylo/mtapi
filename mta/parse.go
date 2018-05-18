package mta

import (
	"os"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/jeffreylo/mtapi/pkg/strings2"
	"github.com/kyroy/kdtree"
	"github.com/kyroy/kdtree/points"
)

// separateStations are IDs of stations that should be considered
// physically separate.
var separateStations = []string{"A27", "132"}

// canonicalID remaps IDs to their canonical station ID.
var canonicalID = map[string]string{
	"635": "L03", // Union Sq - 14 St :: 14 St - Union Sq
	"R20": "L03", // Union Sq - 14 St :: 14 St - Union Sq
	"725": "127", // Times Sq - 42 St
	"902": "127", // Times Sq - 42 St
	"R16": "127", // Times Sq - 42 St
}

// Parser returns station data from GTFS stop and transfer files.
type Parser struct{ StopsPath, TransfersPath string }

// parseResult returns the result of processing.
type parseResult struct {
	StationMap map[string]StationID
	Stations   Stations
	Tree       *kdtree.KDTree
}

// Parse parses the configuration files to create Stations.
func (p *Parser) Parse() (*parseResult, error) {
	type stopRow struct {
		ID            string  `csv:"stop_id"`
		Name          string  `csv:"stop_name"`
		Lat           float64 `csv:"stop_lat"`
		Lon           float64 `csv:"stop_lon"`
		Type          int     `csv:"location_type"`
		ParentStation string  `csv:"parent_station"`
	}

	type transferRow struct {
		FromStopID      string `csv:"from_stop_id"`
		ToStopID        string `csv:"to_stop_id"`
		TransferType    int    `csv:"transfer_type"`
		MinTransferTime int32  `csv:"min_transfer_time"`
	}

	isSeparateStation := func(t *transferRow) bool {
		return !strings2.SliceContains(separateStations, t.FromStopID) && !strings2.SliceContains(separateStations, t.ToStopID)
	}

	s, err := os.OpenFile(p.StopsPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}

	t, err := os.OpenFile(p.TransfersPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}

	var stopRows []*stopRow
	if err := gocsv.UnmarshalFile(s, &stopRows); err != nil {
		panic(err)
	}

	var transferRows []*transferRow
	if err := gocsv.UnmarshalFile(t, &transferRows); err != nil {
		panic(err)
	}

	type stop struct {
		ID          string
		Name        string
		Coordinates struct{ Lat, Lon float64 }
	}

	stops := make(map[string]*stop, len(stopRows))
	for _, v := range stopRows {
		if v.ParentStation == "" {
			stops[v.ID] = &stop{
				ID:   v.ID,
				Name: v.Name,
				Coordinates: struct {
					Lat, Lon float64
				}{v.Lat, v.Lon},
			}
		}
	}

	tree := kdtree.New(nil)
	stationMap := make(map[string]StationID)
	stations := make(Stations, len(transferRows))

	now := time.Now().UTC()

	// Group by destination.
	for _, transfer := range transferRows {
		// If we've already processed the destination, skip it.
		if _, ok := stationMap[transfer.ToStopID]; ok {
			continue
		}

		// Some stop IDs need to be remapped because people
		// think of them as the same station in real life.
		originID := transfer.FromStopID
		if remapID, ok := canonicalID[originID]; ok {
			originID = remapID
		}

		id := StationID(originID)
		// If we haven't stored this in our Stations yet, do it.
		if _, ok := stations[id]; !ok {
			v, ok := stops[string(id)]
			if !ok {
				continue
			}

			// Create the station.
			stations[id] = &Station{
				ID:   id,
				Name: v.Name,
				Coordinates: &Coordinates{
					Lat: v.Coordinates.Lat,
					Lon: v.Coordinates.Lon,
				},
				Arrivals: make(map[Direction][]*Arrival),
				Updated:  &now,
			}

			// A station always maps to itself.
			stationMap[originID] = id
			tree.Insert(points.NewPoint([]float64{v.Coordinates.Lat, v.Coordinates.Lon}, id))
		}

		if isSeparateStation(transfer) {
			stationMap[transfer.ToStopID] = id
		} else {
			// A station always maps to itself.
			stationMap[originID] = id
		}
	}

	return &parseResult{
		StationMap: stationMap,
		Stations:   stations,
		Tree:       tree,
	}, nil
}
