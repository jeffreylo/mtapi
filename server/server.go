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
	client     *mta.Client
	port       int
	dispatcher *jsonrpc.MethodRepository
	ensureSSL  bool
}

// Params defines the server dependencies.
type Params struct {
	Client    *mta.Client
	Port      int
	EnsureSSL bool
}

// New returns a server instance with the specified parameters.
func New(p *Params) *Server {
	mr := jsonrpc.NewMethodRepository()
	must(mr.RegisterMethod("GetSystemStatus", GetSystemStatusHandler{client: p.Client, p: protocol.New()}, nil, GetSystemStatusResult{}))
	must(mr.RegisterMethod("GetStations", GetStationsHandler{client: p.Client, p: protocol.New()}, nil, GetStationsResult{}))
	must(mr.RegisterMethod("GetStation", GetStationHandler{client: p.Client, p: protocol.New()}, GetStationParams{}, GetStationResult{}))
	must(mr.RegisterMethod("GetClosest", GetClosestHandler{client: p.Client, p: protocol.New()}, GetClosestParams{}, GetClosestResult{}))

	return &Server{
		client:     p.Client,
		port:       p.Port,
		dispatcher: mr,
		ensureSSL:  p.EnsureSSL,
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

	var rpcHandler http.Handler
	rpcHandler = s.dispatcher
	if s.ensureSSL {
		rpcHandler = ensureSSL(s.dispatcher)
	}

	m.Handler("GET", "/static/*filepath", http.StripPrefix("/static/", http.FileServer(http.Dir("./client/dist"))))
	m.Handler("GET", "/", http.HandlerFunc(serveTemplate))
	m.Handler("POST", "/rpc", rpcHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      m,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return srv.ListenAndServe()
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("client", "templates", "index.html")
	tmpl, err := template.ParseFiles(lp)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	r.ParseForm()

	if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
