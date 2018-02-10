package server

import (
	"context"

	"github.com/intel-go/fastjson"
	"github.com/jeffreylo/mtapi/mta"
	"github.com/jeffreylo/mtapi/server/protocol"
	"github.com/osamingo/jsonrpc"
)

type GetSystemStatusHandler struct {
	client *mta.Client
	p      *protocol.Protocol
}

type GetSystemStatusResult struct{ Service *mta.Service }

func (h GetSystemStatusHandler) ServeJSONRPC(c context.Context, params *fastjson.RawMessage) (interface{}, *jsonrpc.Error) {
	service, err := h.client.GetServiceStatus()
	if err != nil {
		return nil, &jsonrpc.Error{
			Code:    jsonrpc.ErrorCodeInternal,
			Message: err.Error(),
		}
	}
	return GetSystemStatusResult{Service: service}, nil
}
