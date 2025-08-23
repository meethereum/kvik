package main

import (
	"flag"
	"fmt"

	"example.com/m/v2/config"
	"example.com/m/v2/db"
	"example.com/m/v2/web"

	// "github.com/BurntSushi/toml"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var (
	dbLocation = flag.String("db-location", "", "The path to the boltdb database")
	httpAddr   = flag.String("http-addr", "127.0.0.1:8090", "The http address where API server will run")
	configFile = flag.String("config-file", "sharding.toml", "Config file for static sharding")
	shard      = flag.String("db-shard", "", "the db shard you want to access data from")
)

func parseFlags() {
	flag.Parse()

	if *dbLocation == "" {
		log.Fatalf("Must provide the db location flag")
	}

	if *shard == "" {
		log.Fatalf("Must provide flag mentioning the shard")
	}

}

func main() {

	parseFlags()
	c, err := config.ParseFile(*configFile)
	if err != nil {
		log.Fatal("Error parsing config %q: %v", *configFile, err)
	}
	log.Info("Config file %s successfully parsed ", *configFile)

	//this line is to be removed
	fmt.Printf("%v", c)

	db, close, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatal("NewDatabase(%q): %v", *dbLocation, err)
	}

	defer close()

	srv := web.NewServer(db)

	http.HandleFunc("/get", srv.GetKeyHandler)
	http.HandleFunc("/set", srv.SetKeyHandler)

	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
