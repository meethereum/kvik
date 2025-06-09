package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	bolt "go.etcd.io/bbolt"
)

var(
	dbLocation = flag.String("db-location","","The path to the boltdb database")
	httpAddr = flag.String("http-addr","127.0.0.1:8090","The http address where API server will run")
)

func parseFlags(){
	flag.Parse()

	if *dbLocation == ""{
		log.Fatalf("Must provide the db location flag")
	}
}

func main() {

	parseFlags()

	db, err := bolt.Open(*dbLocation, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/get",func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w,"Called get")
	})

	http.HandleFunc("/set",func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w,"Called set")
	})

	log.Fatal(http.ListenAndServe(*httpAddr,nil))
}