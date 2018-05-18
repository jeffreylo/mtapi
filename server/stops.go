package server

import (
	"context"

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
	stations := h.client.GetStations()
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

	station, err := h.client.GetStation(mta.StationID(p.ID))
	if err != nil {
		return nil, &jsonrpc.Error{
			Code:    jsonrpc.ErrorCodeInvalidParams,
			Message: err.Error(),
		}
	}
	return GetStationResult{Station: h.p.Station(station)}, nil
}

// GetStationResult describes the response of the GetStations RPC.
type GetStationResult struct{ Station interface{} }

// GetClosestHandler returns the nearest stations.
type GetClosestHandler struct {
	client *mta.Client
	p      *protocol.Protocol
}

// GetClosestParams defines the parameters of the GetClosest RPC.
type GetClosestParams struct {
	Lat, Lon    float64
	NumStations int
}

// ServeJSONRPC implements the jsonrpc handler interface.
func (h GetClosestHandler) ServeJSONRPC(c context.Context, params *fastjson.RawMessage) (interface{}, *jsonrpc.Error) {
	var p GetClosestParams
	if err := jsonrpc.Unmarshal(params, &p); err != nil {
		return nil, err
	}
	stations := h.client.GetClosestStations(&mta.Coordinates{Lat: p.Lat, Lon: p.Lon}, p.NumStations)
	vv := make([]*protocol.Station, 0, len(stations))
	for _, v := range stations {
		vv = append(vv, h.p.Station(v))
	}
	return GetClosestResult{Stations: stations}, nil
}

// GetClosestResult is the result.
type GetClosestResult struct{ Stations interface{} }
