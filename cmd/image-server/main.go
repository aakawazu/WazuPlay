package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	imageServer "github.com/aakawazu/WazuPlay/internal/image-server/router"
	"github.com/aakawazu/WazuPlay/pkg/random"
	"github.com/joho/godotenv"
	"github.com/syndtr/goleveldb/leveldb"
)

func initLeveldb() {
	db, err := leveldb.OpenFile("/wazuplay-files/images/leveldb/", nil)
	if err != nil {
		log.Fatal("Error loading leveldb")
	}
	defer db.Close()
	_, err = db.Get([]byte("latest_folder"), nil)
	if err != nil {
		if err != leveldb.ErrNotFound {
			log.Fatal("Error loading latest folder")
		}

		folderName, err := random.GenerateRandomString()
		if err != nil {
			log.Fatal("Error generating randomstring")
		}

		if err := db.Put([]byte("latest_folder"), []byte(folderName), nil); err != nil {
			log.Fatal("Error generating randomstring")
		}
		if err := db.Put([]byte("filesin_latest_folder"), []byte("0"), nil); err != nil {
			log.Fatal("Error generating randomstring")
		}

		err = os.MkdirAll(fmt.Sprintf("/wazuplay-files/images/files/%s", folderName), 0777)
		if err != nil {
			log.Fatal("Error making folder")
		}
	}
}

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := imageServer.NewRouter()
	fmt.Println("hello, Image")
	log.Fatal(http.ListenAndServe(":8080", router))
}
