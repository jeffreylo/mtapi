package main

import (
	"flag"
	"log"

	"github.com/dcowgill/envflag"
	"github.com/jeffreylo/mtapi/mta"
	"github.com/jeffreylo/mtapi/server"
)

func main() {
	var (
		apiKey    = flag.String("api-key", "", "API key from http://datamine.mta.info/")
		ensureSSL = flag.Bool("ensure-ssl", true, "always redirect to https://")
		port      = flag.Int("port", 3000, "port for server")
		path      = flag.String("gtfs-path", "", "gtfs directory")
	)

	flag.Parse()

	envflag.Parse()

	if *apiKey == "" {
		log.Fatal("missing apiKey")
	}
	if *path == "" {
		log.Fatal("missing path")
	}

	client, err := mta.NewClient(&mta.ClientConfig{
		APIKey:            *apiKey,
		StopsFilePath:     *path + "/stops.txt",
		TransfersFilePath: *path + "/transfers.txt",
	})
	if err != nil {
		log.Fatal(err)
	}
	go client.Work()

	server := server.New(&server.Params{
		Client:    client,
		EnsureSSL: *ensureSSL,
		Port:      *port,
	})
	log.Fatal(server.Serve())
}
