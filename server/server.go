package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ubclaunchpad/rocket/config"
	"github.com/ubclaunchpad/rocket/data"
	"github.com/ubclaunchpad/rocket/model"
)

// Server represents the HTTP server that provides a REST API
// interface to Rocket's database.
type Server struct {
	router *mux.Router
	addr   string
	dal    *data.DAL
}

func New(c *config.Config) *Server {
	router := mux.NewRouter()
	s := &Server{
		router: router,
		addr:   c.Host + ":" + c.Port,
	}

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc(RouteMembers, s.MemberHandler).Methods("GET")

	return s
}

func (s *Server) MemberHandler(res http.ResponseWriter, req *http.Request) {
	var members model.Members
	if err := s.dal.GetMembers(&members); err != nil {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(&members)
	}
}
