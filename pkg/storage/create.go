package storage

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aakawazu/WazuPlay/pkg/random"
	"github.com/syndtr/goleveldb/leveldb"
)

// Create create new file
func Create(rootFolderURI string) (*os.File, string, error) {
	db, err := leveldb.OpenFile(fmt.Sprintf("%s/storage/leveldb", rootFolderURI), nil)
	if err != nil {
		return nil, "", err
	}
	defer db.Close()

	fileName, err := random.GenerateRandomString()
	if err != nil {
		return nil, "", err
	}

	folderName, err := db.Get([]byte("latest_folder"), nil)
	if err != nil {
		return nil, "", err
	}

	fileList, err := ioutil.ReadDir(fmt.Sprintf("%s/files/%s", rootFolderURI, folderName))
	if err != nil {
		return nil, "", err
	}

	numberOfFiles := len(fileList)
	if err != nil {
		return nil, "", err
	}

	if numberOfFiles > 50000 {
		folderNameString, err := random.GenerateRandomString()
		if err != nil {
			return nil, "", err
		}

		folderName = []byte(folderNameString)

		err = os.MkdirAll(fmt.Sprintf("%s/files/%s", rootFolderURI, folderName), 0777)
		if err != nil {
			return nil, "", err
		}

		if err := db.Put([]byte("latest_folder"), []byte(folderName), nil); err != nil {
			return nil, "", err
		}
	}

	file, err := os.Create(fmt.Sprintf("%s/files/%s/%s", rootFolderURI, folderName, fileName))
	if err != nil {
		return nil, "", err
	}

	if err := db.Put([]byte(fileName), []byte(folderName), nil); err != nil {
		return nil, "", err
	}

	return file, fileName, nil
}
