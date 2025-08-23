package main

import (
	"flag"

	"example.com/m/v2/config"
	"example.com/m/v2/db"
	"example.com/m/v2/web"

	"net/http"

	log "github.com/sirupsen/logrus"
)

var (
	dbLocation = flag.String("db-location", "", "The path to the boltdb database")
	httpAddr   = flag.String("http-addr", "", "The http address where API server will run")
	configFile = flag.String("config-file", "sharding.toml", "Config file for static sharding")
	shard      = flag.String("db-shard", "", "the db shard you want to access data from")
)

func parseFlags() {
	flag.Parse()

	if *httpAddr == "" {
		log.Fatal("Must provide http address of the db node/shard")
	}

	if *dbLocation == "" {
		log.Fatal("Must provide the db location flag")
	}

	if *shard == "" {
		log.Fatal("Must provide flag mentioning the shard")
	}

}

func main() {

	parseFlags()
	c, err := config.ParseFile(*configFile)
	if err != nil {
		log.Fatalf("Error parsing config %q: %v", *configFile, err)
	}
	log.Infof("Config file %s successfully parsed", *configFile)

	db, close, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatalf("NewDatabase(%q): %v", *dbLocation, err)
	}

	defer close()

	shards := c.Shards
	shardIndex := config.GetShardIndexFromName(*shard, c.Shards)
	addressMappings := config.CreateAddressMappings(shards)
	if shardIndex == -1 {
		log.Fatalf("given shard %s doesnt have an instance. please recheck config for shards", *shard)
	}
	srv := web.NewServer(db, len(shards), shardIndex, addressMappings)

	http.HandleFunc("/get", srv.GetKeyHandler)
	http.HandleFunc("/set", srv.SetKeyHandler)

	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
