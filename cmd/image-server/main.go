package main

import (
	"fmt"
	"log"
	"net/http"

	imageServer "github.com/aakawazu/WazuPlay/internal/image-server/router"
	"github.com/aakawazu/WazuPlay/pkg/upload"
	"github.com/joho/godotenv"
	"github.com/syndtr/goleveldb/leveldb"
)

func initImageFolder() {
	rootFolder := imageServer.ImageFilesRoot
	db, err := leveldb.OpenFile(fmt.Sprintf("%s/leveldb/", rootFolder), nil)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := db.Get([]byte("latest_folder"), nil); err != nil {
		if err != leveldb.ErrNotFound {
			log.Fatal(err)
		}
		if folderName, err := upload.CreateNewFolder(rootFolder); err != nil {
			log.Fatal(err)
		} else {
			if err := db.Put([]byte("latest_folder"), []byte(folderName), nil); err != nil {
				log.Fatal(err)
			}
		}
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
