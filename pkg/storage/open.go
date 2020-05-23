package storage

import (
	"fmt"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
)

// Open open file
func Open(rootFolderURI string, fileName string) (*os.File, error) {
	db, err := leveldb.OpenFile(fmt.Sprintf("%s/storage/leveldb", rootFolderURI), nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	folderName, err := db.Get([]byte(fileName), nil)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(fmt.Sprintf("%s/storage/files/%s/%s", rootFolderURI, folderName, fileName))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return file, nil
}
