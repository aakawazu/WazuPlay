package main

import (
	"fmt"
	"log"
	"net/http"

	oauthServer "github.com/aakawazu/WazuPlay/internal/oauth-server/router"
)

func main() {
	router := oauthServer.NewRouter()
	fmt.Println("http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
