package main

import (
	"flag"
	"log"

	"github.com/jeffreylo/mtapi/mta"
)

func main() {
	var (
		apiKey string
		port   int
		path   string
	)

	flag.StringVar(&apiKey, "apiKey", "", "http://datamine.mta.info/")
	flag.StringVar(&path, "path", "", "path to gtfs stops.txt")
	flag.IntVar(&port, "port", 3000, "port for server")

	flag.Parse()

	if apiKey == "" {
		log.Fatal("missing apiKey")
	}
	if path == "" {
		log.Fatal("missing path")
	}

	client, err := mta.NewClient(&mta.ClientConfig{
		APIKey:        apiKey,
		Port:          port,
		StopsFilePath: path,
	})
	if err != nil {
		log.Fatal(err)
	}
	go client.Work()
	log.Fatal(client.Serve())
}
