package main

import (
	"example.com/m/v2/db"
	"example.com/m/v2/web"
	"flag"
	"log"
	"net/http"
)

var (
	dbLocation = flag.String("db-location", "", "The path to the boltdb database")
	httpAddr   = flag.String("http-addr", "127.0.0.1:8090", "The http address where API server will run")
)

func parseFlags() {
	flag.Parse()

	if *dbLocation == "" {
		log.Fatalf("Must provide the db location flag")
	}
}

func main() {

	parseFlags()

	db, close, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatalf("NewDatabase(%q): %v", *dbLocation, err)
	}

	defer close()

	srv := web.NewServer(db)

	http.HandleFunc("/get", srv.GetKeyHandler)
	http.HandleFunc("/set", srv.SetKeyHandler)

	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
