package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jeffreylo/mtapi/mta"
	"github.com/jeffreylo/mtapi/server/protocol"
	"github.com/osamingo/jsonrpc"
)

// Server is the interface to the MTA API.
type Server struct {
	client     *mta.Client
	port       int
	dispatcher *jsonrpc.MethodRepository
}

// Params defines the server dependencies.
type Params struct {
	Client *mta.Client
	Port   int
}

// New returns a server instance with the specified parameters.
func New(p *Params) *Server {
	mr := jsonrpc.NewMethodRepository()
	must(mr.RegisterMethod("GetSystemStatus", GetSystemStatusHandler{client: p.Client, p: protocol.New()}, nil, GetSystemStatusResult{}))
	must(mr.RegisterMethod("GetStations", GetStationsHandler{client: p.Client, p: protocol.New()}, nil, GetStationsResult{}))
	must(mr.RegisterMethod("GetStation", GetStationHandler{client: p.Client, p: protocol.New()}, GetStationParams{}, GetStationResult{}))

	return &Server{
		client:     p.Client,
		port:       p.Port,
		dispatcher: mr,
	}
}

// Serve returns an http server.
func (s *Server) Serve() error {
	m := http.NewServeMux()
	m.Handle("/rpc", s.dispatcher)
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      m,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return srv.ListenAndServe()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
