package gateway

import (
	"net/http"
	"github.com/pressly/chi"
	"github.com/gocql/gocql"
	"github.com/pressly/chi/render"

	"context"
	"zeats/config"
	"fmt"
	"github.com/golang/glog"
	"zeats/worker"
	"net"
)

type Server struct {
	httpServer *http.Server
	router *chi.Mux
}

// listenAndServe would stast listening on given port
func (s *Server) listenAndServe() error {

	listener, err := net.Listen("tcp", ":"+s.httpServer.Addr)
	if err != nil {
		return err
	}

	return s.httpServer.Serve(listener)
}

// newServer would init server
func newServer(cSession *gocql.Session) *Server {
	r := chi.NewRouter()

	r.Use(render.SetContentType(render.ContentTypeJSON))

	zeatsService := newService(&Config{CassandraSession: cSession})
	setupRoutes(r, zeatsService)

	server := &Server{
		httpServer: &http.Server{Addr: config.Vals().Port, Handler: r},
		router:     r,
	}

	return server
}

func setupRoutes(r *chi.Mux, zeatsService Service) {
	r.Route("/", func(r chi.Router) {
		r.Use(apiVersionCtx("v1"))
		r.Mount("/v1/", Handler(zeatsService))
	})
}

func apiVersionCtx(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), "api.version", version))
			next.ServeHTTP(w, r)
		})
	}
}

type Gateway interface {
	Produce()
	Start()
	Stop()
	Pool() *worker.SourceWorkerPool
}

type gateway struct {
	s    *Server
	c    *gocql.Session
	pool *worker.SourceWorkerPool
}

func (g *gateway) Produce() {
	g.s = newServer(g.c)
	g.s.listenAndServe()
}

func (g *gateway) Start() {
	fmt.Println("[✔]\tGateway started")
	g.pool.Start()
	g.pool.Wait()
}

func (g *gateway) Stop() {
	if err := g.s.httpServer.Shutdown(context.Background()); err != nil {
		glog.Fatalf("Could not shutdown: %v", err)
	}
	g.pool.Stop()
	g.c.Close()
	fmt.Println("[✘]\tGateway stopped")
}

func (g *gateway) Pool() *worker.SourceWorkerPool {
	return g.pool
}

func NewGateway(c *gocql.Session) Gateway {
	gw := &gateway{c: c}
	gw.pool = worker.NewSourceWorkerPool("gateway", 1, gw)
	return gw
}


