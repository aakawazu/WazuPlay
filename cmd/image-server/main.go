package main

import (
	"fmt"
	"log"
	"net/http"

	imageServer "github.com/aakawazu/WazuPlay/internal/image-server/router"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	router := imageServer.NewRouter()
	fmt.Println("hello, Image")
	log.Fatal(http.ListenAndServe(":8080", router))
}
