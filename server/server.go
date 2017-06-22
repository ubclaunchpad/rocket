package server

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

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

// New returns a new instance of the HTTP server based on a config.
func New(c *config.Config, dal *data.DAL) *Server {
	router := mux.NewRouter()
	s := &Server{
		router: router,
		addr:   c.Host + ":" + c.Port,
		dal:    dal,
	}

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/members", s.MemberHandler).Methods("GET")
	api.HandleFunc("/teams", s.TeamHandler).Methods("GET")

	return s
}

func (s *Server) Start() error {
	return http.ListenAndServe(s.addr, s.router)
}

func (s *Server) MemberHandler(res http.ResponseWriter, req *http.Request) {
	var members model.Members
	if err := s.dal.GetMembers(&members); err != nil {
		log.WithError(err).Error("Failed to get members")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(res).Encode(&members); err != nil {
		log.WithError(err).Error("Failed to encode JSON")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func (s *Server) TeamHandler(res http.ResponseWriter, req *http.Request) {
	var teams model.Teams
	if err := s.dal.GetTeams(&teams); err != nil {
		log.WithError(err).Error("Failed to get teams")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(res).Encode(&teams); err != nil {
		log.WithError(err).Error("Failed to encode JSON")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}
