package main

import (
	"flag"
	"log"

	"github.com/cfmobile/gopivnet/api"
	"github.com/cfmobile/tile-migrations-generator/migrations"
)

func main() {
	var product, token, path string
	flag.StringVar(&product, "p", "", "product slug name")
	flag.StringVar(&token, "t", "", "pivnet token")
	flag.StringVar(&path, "o", "", "content migrations path")

	flag.Parse()
	if product == "" || token == "" || path == "" {
		log.Fatalf("Invalid flags\n")
	}

	api := api.New(token)
	fetcher := migrations.NewProductVersionFetcher()

	contentMigrations, err := migrations.New(api, fetcher, path)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = contentMigrations.WriteMissingMigrations(product)
	if err != nil {
		log.Fatal(err)
	}
}
