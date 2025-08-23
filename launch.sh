#!/bin/bash
set -e

trap 'killall kvik' SIGINT

cd "$(dirname "$0")"

# Kill old kvik processes
killall kvik 2>/dev/null || true
sleep 0.1

# Build the binary locally
go build -o kvik main.go

# Start shard instances
./kvik -db-location=mumbai.db    -http-addr=127.0.0.2:8080 -config-file=sharding.toml -db-shard=Mumbai &
./kvik -db-location=singapore.db -http-addr=127.0.0.3:8080 -config-file=sharding.toml -db-shard=Singapore &
./kvik -db-location=newyork.db   -http-addr=127.0.0.4:8080 -config-file=sharding.toml -db-shard=NewYork &
./kvik -db-location=moscow.db    -http-addr=127.0.0.5:8080 -config-file=sharding.toml -db-shard=Moscow &

wait
