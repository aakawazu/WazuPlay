package main

import (
	"fmt"
	"log"
	"net/http"

	audioServer "github.com/aakawazu/WazuPlay/internal/audio-server/router"
	"github.com/aakawazu/WazuPlay/pkg/storage"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	storage.Init(audioServer.AudioFilesRoot)

	router := audioServer.NewRouter()
	fmt.Println("hello, Audio")
	log.Fatal(http.ListenAndServe(":8080", router))
}
