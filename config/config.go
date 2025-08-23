package config

import (
	"hash/fnv"

	"github.com/BurntSushi/toml"
)

type Shard struct {
	Name    string
	Idx     int
	Address string
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

func GetShardIndexFromName(shardName string, shards []Shard) int {
	var i int = -1
	for i = range shards {
		if shards[i].Name == shardName {
			break
		}
	}
	if i != -1 {
		return shards[i].Idx
	}
	return -1
}

func CreateAddressMappings(shards []Shard) map[int]string {
	mappings := make(map[int]string)
	for i := range shards {
		mappings[shards[i].Idx] = shards[i].Address
	}
	return mappings
}
