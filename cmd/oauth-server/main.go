package main

import (
	"log"
	"net/http"

	oauthServer "github.com/aakawazu/WazuPlay/internal/oauth-server/router"
)

func main() {
	router := oauthServer.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
