package server

import (
	"context"
	"fmt"
	"time"

	"github.com/intel-go/fastjson"
	"github.com/jeffreylo/mtapi/mta"
	"github.com/jeffreylo/mtapi/server/protocol"
	"github.com/osamingo/jsonrpc"
)

// GetStationsHandler returns all stations.
type GetStationsHandler struct {
	client *mta.Client
	p      *protocol.Protocol
}

// ServeJSONRPC implements the jsonrpc handler interface.
func (h GetStationsHandler) ServeJSONRPC(c context.Context, params *fastjson.RawMessage) (interface{}, *jsonrpc.Error) {
	stations := h.client.Stations()
	return GetStationsResult{Stations: h.p.Stations(stations)}, nil
}

// GetStationsResult describes the response of the GetStations RPC.
type GetStationsResult struct{ Stations []*protocol.Station }

// GetStationHandler returns all schedules for a given station.
type GetStationHandler struct {
	client *mta.Client
	p      *protocol.Protocol
}

// GetStationParams defines the parameters of the GetStation RPC.
type GetStationParams struct{ ID string }

// ServeJSONRPC implements the jsonrpc handler interface.
func (h GetStationHandler) ServeJSONRPC(c context.Context, params *fastjson.RawMessage) (interface{}, *jsonrpc.Error) {
	var p GetStationParams
	if err := jsonrpc.Unmarshal(params, &p); err != nil {
		return nil, err
	}

	stops := h.client.GetStops()
	stations := h.client.Stations()
	station, ok := stations[mta.StopID(p.ID)]
	if !ok {
		return nil, &jsonrpc.Error{
			Code:    jsonrpc.ErrorCodeInvalidParams,
			Message: fmt.Sprintf("Station ID=%s is invalid or does not exist.", p.ID),
		}
	}

	stationSchedule := make(map[mta.Direction]mta.Schedule)
	var updated *time.Time
	for _, id := range station.StopIDs() {
		stop := stops[id]
		for d, s := range stop.Schedules {
			if _, ok := stationSchedule[d]; !ok {
				stationSchedule[d] = make(mta.Schedule, 0, len(s))
			}
			stationSchedule[d] = append(stationSchedule[d], s...)
		}
		updated = stop.Updated
	}

	return GetStationResult{Station: h.p.Station(station, stationSchedule, updated)}, nil
}

// GetStationResult describes the response of the GetStations RPC.
type GetStationResult struct{ Station interface{} }
