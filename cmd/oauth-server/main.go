package main

import (
	"fmt"
	"log"
	"net/http"

	oauthServer "github.com/aakawazu/WazuPlay/internal/oauth-server/router"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	router := oauthServer.NewRouter()
	fmt.Println("hello, OAuth")
	log.Fatal(http.ListenAndServe(":8080", router))
}
