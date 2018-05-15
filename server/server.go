package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/jeffreylo/mtapi/mta"
	"github.com/jeffreylo/mtapi/server/protocol"
	"github.com/julienschmidt/httprouter"
	"github.com/osamingo/jsonrpc"
)

// Server is the interface to the MTA API.
type Server struct {
	client      *mta.Client
	dispatcher  *jsonrpc.MethodRepository
	ensureSSL   bool
	environment string
	port        int
	release     string
	staticPath  string
}

// Params defines the server dependencies.
type Params struct {
	Client      *mta.Client
	EnsureSSL   bool
	Environment string
	Port        int
	Release     string
	StaticPath  string
}

// New returns a server instance with the specified parameters.
func New(p *Params) *Server {
	mr := jsonrpc.NewMethodRepository()
	must(mr.RegisterMethod("GetSystemStatus", GetSystemStatusHandler{client: p.Client, p: protocol.New()}, nil, GetSystemStatusResult{}))
	must(mr.RegisterMethod("GetStations", GetStationsHandler{client: p.Client, p: protocol.New()}, nil, GetStationsResult{}))
	must(mr.RegisterMethod("GetStation", GetStationHandler{client: p.Client, p: protocol.New()}, GetStationParams{}, GetStationResult{}))
	must(mr.RegisterMethod("GetClosest", GetClosestHandler{client: p.Client, p: protocol.New()}, GetClosestParams{}, GetClosestResult{}))

	return &Server{
		client:      p.Client,
		dispatcher:  mr,
		ensureSSL:   p.EnsureSSL,
		environment: p.Environment,
		port:        p.Port,
		release:     p.Release,
		staticPath:  p.StaticPath,
	}
}

func ensureSSL(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Forwarded-Proto") != "https" {
			target := "https://" + r.Host + r.URL.Path
			if len(r.URL.RawQuery) > 0 {
				target += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, target, http.StatusTemporaryRedirect)
		}
		next.ServeHTTP(w, r)
	})
}

// Serve returns an http server.
func (s *Server) Serve() error {
	m := httprouter.New()

	var fileHandler, indexHandler, rpcHandler http.Handler
	fileHandler = http.StripPrefix("/static/", http.FileServer(http.Dir(s.staticPath)))
	indexHandler = serveTemplate(&tmplData{s.environment, s.release})
	rpcHandler = s.dispatcher
	if s.ensureSSL {
		fileHandler = ensureSSL(fileHandler)
		indexHandler = ensureSSL(indexHandler)
		rpcHandler = ensureSSL(rpcHandler)
	}

	m.Handler("GET", "/static/*filepath", fileHandler)
	m.Handler("GET", "/", indexHandler)
	m.Handler("POST", "/rpc", rpcHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      m,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return srv.ListenAndServe()
}

type tmplData struct{ Environment, Release string }

func serveTemplate(data *tmplData) http.HandlerFunc {
	lp := filepath.Join("client", "templates", "index.html")
	tmpl, err := template.ParseFiles(lp)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
			return
		}
		if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
	})
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
