package storage

import (
	"fmt"
	"log"
	"os"

	"github.com/aakawazu/WazuPlay/pkg/random"
	"github.com/syndtr/goleveldb/leveldb"
)

func Init(rootFolderUri string) error {
	db, err := leveldb.OpenFile(fmt.Sprintf("%s/leveldb", rootFolderUri), nil)
	defer db.Close()
	if err != nil {
		return err
	}

	if _, err := db.Get([]byte("latest_folder"), nil); err != nil {
		if err != leveldb.ErrNotFound {
			log.Fatal(err)
		}

		folderName, err := random.GenerateRandomString()
		if err != nil {
			return err
		}

		err = os.MkdirAll(fmt.Sprintf("%s/files/%s", rootFolderUri, folderName), 0777)
		if err != nil {
			return err
		}

		if err := db.Put([]byte("latest_folder"), []byte(folderName), nil); err != nil {
			return err
		}
	}

	return nil
}
