package mta

import (
	"os"
	"strings"

	"github.com/Kyroy/kdtree"
	"github.com/Kyroy/kdtree/points"
	"github.com/gocarina/gocsv"
	"github.com/jeffreylo/mtapi/pkg/strings2"
)

// Parse parses the configuration files to create a Stops and
// Stations.
func Parse(stopsPath string, transfersPath string) (Stops, Stations, *kdtree.KDTree, error) {
	type stopRow struct {
		ID            StopID  `csv:"stop_id"`
		Name          string  `csv:"stop_name"`
		Lat           float64 `csv:"stop_lat"`
		Lon           float64 `csv:"stop_lon"`
		Type          int     `csv:"location_type"`
		ParentStation string  `csv:"parent_station"`
	}

	type transferRow struct {
		FromStopID      StopID `csv:"from_stop_id"`
		ToStopID        StopID `csv:"to_stop_id"`
		TransferType    int    `csv:"transfer_type"`
		MinTransferTime int32  `csv:"min_transfer_time"`
	}

	isSeparateStation := func(t *transferRow) bool {
		var separateStations = []string{"A27", "132"}
		return !strings2.SliceContains(separateStations, string(t.FromStopID)) && !strings2.SliceContains(separateStations, string(t.ToStopID))
	}

	nameMap := map[StopID]StopID{"L03": "R20"}

	if stopsPath == "" {
		stopsPath = defaultStopsFile
	}

	if transfersPath == "" {
		transfersPath = defaultTransfersFile
	}

	s, err := os.OpenFile(stopsPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}

	t, err := os.OpenFile(transfersPath, os.O_RDONLY, os.ModePerm)
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

	stops := make(Stops, len(stopRows))
	for _, stop := range stopRows {
		if stop.ParentStation == "" {
			stops[stop.ID] = &Stop{
				ID:   stop.ID,
				Name: stop.Name,
				Coordinates: &Coordinates{
					Lat: stop.Lat,
					Lon: stop.Lon,
				},
				Schedules: make(map[Direction]Schedule),
			}
		}
	}

	sentinel := make(map[StopID]struct{})
	tree := kdtree.New(nil)
	stations := make(Stations, len(transferRows))
	for _, transfer := range transferRows {
		if _, ok := sentinel[transfer.ToStopID]; ok {
			continue
		}

		_, originExists := stations[transfer.FromStopID]
		if !originExists {
			v, ok := stops[transfer.FromStopID]
			if !ok {
				continue
			}

			stations[transfer.FromStopID] = &Station{
				ID:          transfer.FromStopID,
				Coordinates: v.Coordinates,
				Stops:       make(map[StopID]struct{}),
			}
			stations[transfer.FromStopID].Stops[transfer.FromStopID] = struct{}{}
			tree.Insert(points.NewPoint([]float64{v.Coordinates.Lat, v.Coordinates.Lon}, transfer.FromStopID))
			sentinel[transfer.FromStopID] = struct{}{}
		}

		if isSeparateStation(transfer) {
			stations[transfer.FromStopID].Stops[transfer.ToStopID] = struct{}{}
			sentinel[transfer.ToStopID] = struct{}{}
		} else {
			stations[transfer.FromStopID].Stops[transfer.FromStopID] = struct{}{}
			sentinel[transfer.FromStopID] = struct{}{}
		}
	}

	for _, station := range stations {
		names := make([]string, 0, len(station.StopIDs()))
		for _, id := range station.StopIDs() {
			if remapID, ok := nameMap[id]; ok {
				id = remapID
			}
			names = append(names, stops[id].Name)
		}
		names = strings2.Unique(names)
		station.Name = strings.Join(names, " / ")
	}

	return stops, stations, tree, nil
}
