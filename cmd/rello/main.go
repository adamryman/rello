package main

import (
	"net/http"
	"os"

	"github.com/adamryman/rello"
)

func main() {
	port := os.Getenv("PORT")
	dbLocation := os.Getenv("SQLITE3")

	err := rello.InitDatabase(dbLocation)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", rello.HandleWebhook)
	_ = http.ListenAndServe(":"+port, nil)
}
