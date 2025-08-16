package web

import (
	"fmt"
	"net/http"

	"example.com/m/v2/db"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	db *db.Database
}

func NewServer(db *db.Database) *Server {
	return &Server{db: db}
}

func (s *Server) GetKeyHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value, err := s.db.GetKey(key)
	fmt.Fprintf(w, "Value = %q, error = %v", value, err)

	log.WithFields(log.Fields{
		"key":   key,
		"value": value,
		"error": err,
	}).Info("GetKeyHandler invoked")
}

func (s *Server) SetKeyHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Info("SetKeyHandler has been invoked")

	key := r.Form.Get("key")
	value := r.Form.Get("value")

	log.WithFields(log.Fields{
		"key":   key,
		"value": value,
	}).Info("user requested key be set")

	err := s.db.SetKey(key, []byte(value))
	fmt.Fprintf(w, "Value = %q, error = %v", value, err)

	log.WithFields(log.Fields{
		"key":   key,
		"value": value,
		"error": err,
	}).Info("SetKeyHandler completed")
}
