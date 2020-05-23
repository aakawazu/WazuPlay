package storage

import (
	"fmt"
	"log"
	"os"

	"github.com/aakawazu/WazuPlay/pkg/random"
	"github.com/syndtr/goleveldb/leveldb"
)

// Init init file storage
func Init(rootFolderURI string) error {
	db, err := leveldb.OpenFile(fmt.Sprintf("%s/leveldb", rootFolderURI), nil)
	if err != nil {
		return err
	}
	defer db.Close()

	if _, err := db.Get([]byte("latest_folder"), nil); err != nil {
		if err != leveldb.ErrNotFound {
			log.Fatal(err)
		}

		folderName, err := random.GenerateRandomString()
		if err != nil {
			return err
		}

		err = os.MkdirAll(fmt.Sprintf("%s/files/%s", rootFolderURI, folderName), 0777)
		if err != nil {
			return err
		}

		if err := db.Put([]byte("latest_folder"), []byte(folderName), nil); err != nil {
			return err
		}
	}

	return nil
}
