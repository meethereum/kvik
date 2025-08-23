package web

import (
	"example.com/m/v2/config"
	"example.com/m/v2/db"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type Server struct {
	db              *db.Database
	shardCount      int
	shardIndex      int            //index of shard/server the db is pointing to
	addressMappings map[int]string //mappings of shard index and their corresponding locations
}

func NewServer(db *db.Database, shardCount int, shardIndex int, addressMappings map[int]string) *Server {
	return &Server{db: db, shardCount: shardCount, shardIndex: shardIndex, addressMappings: addressMappings}
}

func (s *Server) GetKeyHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")

	shardNumber, err := config.FindShardIndex([]byte(key), s.shardCount)

	if err != nil {
		log.Error("Error while finding shard index %v", s)
	}

	log.WithFields(log.Fields{
		"key": key,
		"shardNumber":shardNumber,
	}).Info("user requested to get the value")

	if shardNumber == s.shardIndex {
		value, err := s.db.GetKey(string(key))
		fmt.Fprintf(w, "Value = %q, error = %v", value, err)
	} else {
		correctAddress := s.addressMappings[shardNumber]
		// here we dont need to worry about pass by value since w is an interface and its reference will eb only passed
		redirectGetRequest(correctAddress, []byte(key), w)
	}

}

func (s *Server) SetKeyHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Info("SetKeyHandler has been invoked")

	key := r.Form.Get("key")
	value := r.Form.Get("value")

	shardNumber, err := config.FindShardIndex([]byte(key), s.shardCount)

	if err != nil {
		log.Error("Error while finding shard index %v", s)
	}

	log.WithFields(log.Fields{
		"key":   key,
		"value": value,
		"shardNumber":shardNumber,
	}).Info("user requested key be set")

	if shardNumber == s.shardIndex {
		err = s.db.SetKey(key, []byte(value))
		fmt.Fprintf(w, "Value = %q, error = %v", value, err)
	} else {
		correctAddress := s.addressMappings[shardNumber]
		// here we dont need to worry about pass by value since w is an interface and its reference will eb only passed
		redirectSetRequest(correctAddress, []byte(key), []byte(value), w)
	}

	log.WithFields(log.Fields{
		"key":   key,
		"value": value,
		"error": err,
	}).Info("SetKeyHandler completed successfully")
}

func redirectSetRequest(correctAddress string, key []byte, value []byte, w http.ResponseWriter) {

	forwardURL := fmt.Sprintf("http://%s/set?key=%s&value=%s", correctAddress, key, value)

	resp, err := http.Get(forwardURL)
	if err != nil {
		http.Error(w, "Forwarding failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	//todo : check for internal server(shard error)
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

func redirectGetRequest(correctAddress string, key []byte, w http.ResponseWriter) {
	forwardURL := fmt.Sprintf("http://%s/get?key=%s", correctAddress, key)

	resp, err := http.Get(forwardURL)
	if err != nil {
		http.Error(w, "Forwarding failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Failed to copy response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
