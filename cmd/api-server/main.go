package main

import (
	"fmt"
	"log"
	"net/http"

	apiServer "github.com/aakawazu/WazuPlay/internal/api-server/router"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := apiServer.NewRouter()
	fmt.Println("hello, API")
	log.Fatal(http.ListenAndServe(":8080", router))
}
