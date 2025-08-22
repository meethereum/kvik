package config

import (
	"github.com/BurntSushi/toml"
	"hash/fnv"
)

type Shard struct {
	Name string
	Idx  int
}

type Config struct {
	Shards []Shard
}

// ParseFile takes a filename as input and parses it and returns a config
func ParseFile(filename string) (Config, error) {
	var c Config
	if _, err := toml.DecodeFile(filename, &c); err != nil {
		return Config{}, err
	}
	return c, nil
}




// FindShardIndex returns the shard number for corresponding key
func FindShardIndex(key []byte, numShards int) (int, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(key))
	if err != nil {
		return -1, err
	}
	return int(h.Sum32()) % numShards, nil
}
