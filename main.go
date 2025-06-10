package main

import (
	"flag"
	"log"
	"net/http"

	"example.com/m/v2/config"
	"example.com/m/v2/db"
	"example.com/m/v2/web"

	"github.com/BurntSushi/toml"
)

var (
	dbLocation = flag.String("db-location", "", "The path to the boltdb database")
	httpAddr   = flag.String("http-addr", "127.0.0.1:8090", "The HTTP address where API server will run")
	configFile = flag.String("config-file", "sharding.toml", "Path to the sharding config file")
	shard      = flag.String("shard", "", "The name of the shard for the data")
)

func parseFlags() config.Config {
	flag.Parse()

	if *dbLocation == "" {
		log.Fatalf("Must provide the db location flag")
	}

	var c config.Config
	_, err := toml.DecodeFile(*configFile, &c)
	if err != nil {
		log.Fatalf("toml.DecodeFile(%q): %v", *configFile, err)
	}

	return c
}

func main() {
	c := parseFlags()

	var shardCount int
	shardIdx := -1

	for _, s := range c.Shards {
		if s.Idx+1 > shardCount {
			shardCount = s.Idx + 1
		}
		if s.Name == *shard {
			shardIdx = s.Idx
		}
	}

	if shardIdx < 0 {
		log.Fatalf("Shard %q not found", *shard)
	}

	log.Printf("Shard count is %d, current shard is %d", shardCount, shardIdx)

	database, closeFunc, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatalf("NewDatabase(%q): %v", *dbLocation, err)
	}
	defer closeFunc()

	srv := web.NewServer(database)

	http.HandleFunc("/get", srv.GetKeyHandler)
	http.HandleFunc("/set", srv.SetKeyHandler)

	log.Printf("Server starting at %s...", *httpAddr)
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
