package main

import (
	"fmt"
	"log"
	"net/http"

	imageServer "github.com/aakawazu/WazuPlay/internal/image-server/router"
	"github.com/joho/godotenv"
	"github.com/syndtr/goleveldb/leveldb"
)

func initImageFolder() {
	db, err := leveldb.OpenFile(fmt.Sprintf("%s/leveldb/", imageServer.ImageFilesRoot), nil)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := db.Get([]byte("latest_folder"), nil); err != nil {
		if err != leveldb.ErrNotFound {
			log.Fatal(err)
		}
		imageServer.CreateNewFolder(db)
	}
}

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	initImageFolder()

	router := imageServer.NewRouter()
	fmt.Println("hello, Image")
	log.Fatal(http.ListenAndServe(":8080", router))
}
