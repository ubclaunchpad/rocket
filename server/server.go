package server

import (
	"crypto/tls"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/acme/autocert"

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
	server *http.Server
	addr   string
	dal    *data.DAL
	log    *log.Entry
}

// New returns a new instance of the HTTP server based on a config.
func New(c *config.Config, dal *data.DAL, entry *log.Entry) *Server {
	router := mux.NewRouter()
	certManager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(c.Domain),
		Cache:      autocert.DirCache(c.CertificateDir),
	}
	server := &http.Server{
		Addr: ":443",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	s := &Server{
		router: router,
		server: server,
		addr:   c.Host + ":" + c.Port,
		dal:    dal,
		log:    entry,
	}

	router.HandleFunc("/", s.RootHandler).Methods("GET")

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/members", s.MemberHandler).Methods("GET")
	api.HandleFunc("/teams", s.TeamHandler).Methods("GET")

	return s
}

func (s *Server) Start() error {
	s.log.Info("Starting API server on ", s.addr)
	err := s.server.ListenAndServeTLS("", "")
	if err != nil {
		s.log.WithError(err).Error("Error serving HTTP")
	}
	return err
}

func (s *Server) RootHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(`
	<html>
		<head>
		</head>
		<body style="display: flex; align-items: center; justify-content: center; font-size: 64px;">
			&#x1F680;
		</body>
	</html>
	`))
}

func (s *Server) MemberHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	var members model.Members
	if err := s.dal.GetMembers(&members); err != nil {
		s.log.WithError(err).Error("Failed to get members")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(res).Encode(&members); err != nil {
		s.log.WithError(err).Error("Failed to encode JSON")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func (s *Server) TeamHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	var teams model.Teams
	if err := s.dal.GetTeams(&teams); err != nil {
		s.log.WithError(err).Error("Failed to get teams")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(res).Encode(&teams); err != nil {
		s.log.WithError(err).Error("Failed to encode JSON")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}
