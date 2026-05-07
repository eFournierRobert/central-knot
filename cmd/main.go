package main

import (
	"central-knot/internal/api"
	"central-knot/internal/repository/db_utils"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting Central Knot tracker...")

	log.Println("Opening connection to database...")
	if err := db_utils.OpenDatabaseConnection(); err != nil {
		log.Fatal("couldn't connect to database")
	}
	log.Println("Connected to database successfully")

	api.SetUpAnnounceEndpoint()

	println(" +-+-+-+-+-+-+-+ +-+-+-+-+\n |C|e|n|t|r|a|l| |K|n|o|t|\n +-+-+-+-+-+-+-+ +-+-+-+-+")
	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatalln(err)
	}
}
