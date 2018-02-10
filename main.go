package main

import (
	"flag"
	"log"

	"github.com/jeffreylo/mtapi/mta"
	"github.com/jeffreylo/mtapi/server"
)

func main() {
	var (
		apiKey string
		port   int
		path   string
	)

	flag.StringVar(&apiKey, "apiKey", "", "http://datamine.mta.info/")
	flag.StringVar(&path, "gtfs", "", "gtfs directory")
	flag.IntVar(&port, "port", 3000, "port for server")
	flag.Parse()

	if apiKey == "" {
		log.Fatal("missing apiKey")
	}
	if path == "" {
		log.Fatal("missing path")
	}

	client, err := mta.NewClient(&mta.ClientConfig{
		APIKey:            apiKey,
		StopsFilePath:     path + "/stops.txt",
		TransfersFilePath: path + "/transfers.txt",
	})
	if err != nil {
		log.Fatal(err)
	}
	go client.Work()

	server := server.New(&server.Params{
		Client: client,
		Port:   port,
	})
	log.Fatal(server.Serve())
}
