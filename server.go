package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hyperionii/medivac/config"
	"github.com/hyperionii/medivac/handlers"
	"github.com/hyperionii/medivac/handlers/httputils"
	"github.com/hyperionii/medivac/routes"
)

// Server contains instance configuration variables for the server.
type Server struct {
	cfg    *config.Config
	router *mux.Router
}

// NewServer returns a new instance of the server.
func NewServer() (*Server, error) {
	log.Println("Configuring server..")
	cfg, err := config.NewConfig()
	cfg.Print()
	if err != nil {
		return nil, err
	}

	server := &Server{
		cfg: cfg,
	}

	server.configureRouter()
	return server, nil
}

// ListenAndServe attaches the current server to the specified configuration port.
func (s *Server) ListenAndServe() error {
	port := ":" + s.cfg.App.Port

	return http.ListenAndServe(port, s.router)
}

func (s *Server) bindRoutes(r []routes.Route) {
	for _, route := range r {
		handler := s.makeHTTPHandler(route)

		s.router.
			Methods(route.Method()).
			Path(route.Pattern()).
			HandlerFunc(handler)
	}
}

func (s *Server) configureRouter() {
	s.router = mux.NewRouter().StrictSlash(true)
	r := routes.NewRoutes(s.cfg)

	s.bindRoutes(r.APIRoutes)
}

// makeHTTPHandler creates a http.HandlerFunc from an httputils.APIHandler.
func (s *Server) makeHTTPHandler(route routes.Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerFunc := s.handleWithMiddlewares(route)
		err := handlerFunc(w, r)

		if err != nil {
			log.Printf("Handler [%s][%s] returned error: %s", r.Method, r.URL.Path, err)
		}
	}
}

// handleWithMiddlewares applies all middlewares to the specified route. Some
// middleware functions are applied depending on the route's properties, such
// as ValidateAuth and Authorize middlewares. These last 2 functions require
// that the route RequiresAuth() and that RequiredRoles() > 0.
func (s *Server) handleWithMiddlewares(route routes.Route) httputils.APIHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		h := route.HandlerFunc()
		h = handlers.HandleHTTPError(h)
		h = handlers.GzipContent(h)

		if route.RequiresAuth() {
			if requiredRoles := route.RequiredRoles(); len(requiredRoles) > 0 {
				h = handlers.Authorize(h)
			}

			h = handlers.ValidateAuth(h)
		}

		return h(w, r)
	}
}
