package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cfmobile/tile-migrations-generator/migrations"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Not enough args")
	}

	fileName := os.Args[1]

	productVersionFetcher := migrations.NewProductVersionFetcher()
	version, err := productVersionFetcher.FetchProductVersion(fileName)
	if err != nil {
		log.Fatal("Unable to get version. Error: " + err.Error())
	}
	fmt.Println(version)
}
