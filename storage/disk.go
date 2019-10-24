package storage

import (
	"github.com/ad-sho-loko/bogodb/meta"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

// diskManager is responsible for fetching pages from the disk.
type diskManager struct {
}

func newDiskManager() *diskManager {
	return &diskManager{}
}

func (d *diskManager) toPid(tid uint64) uint64 {
	return tid / TupleNumber
}

func (d *diskManager) fetchPage(dirPath, tableName string, pgid uint64) (*Page, error) {
	fileName := strconv.FormatUint(pgid, 10)
	pagePath := path.Join(dirPath, tableName, fileName)

	_, err := os.Stat(pagePath)
	if os.IsNotExist(err) {
		return nil, err
	}

	bytes, err := ioutil.ReadFile(pagePath)
	if err != nil {
		return nil, err
	}

	var b [4096]byte
	copy(b[:], bytes)
	return DeserializePage(b)
}

func (d *diskManager) persist(dirName string, tableName string, pgid uint64, page *Page) error {
	b, err := SerializePage(page)

	if err != nil {
		return err
	}

	p := filepath.Join(dirName, tableName)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		err := os.Mkdir(p, 0777)
		if err != nil {
			panic(err)
		}
	}

	fileName := strconv.FormatUint(pgid, 10)
	savePath := path.Join(dirName, tableName, fileName)
	return ioutil.WriteFile(savePath, b[:], 0644)
}

func (d *diskManager) readIndex(indexName string) (*meta.BTree, error) {
	readPath := path.Join(indexName)
	bytes, err := ioutil.ReadFile(readPath)
	if err != nil {
		return nil, err
	}

	btree, err := meta.DeserializeBTree(bytes)
	if err != nil {
		return nil, err
	}

	return btree, nil
}

func (d *diskManager) writeIndex(dirPath, indexName string, tree *meta.BTree) error {
	b, err := meta.SerializeBTree(tree)

	if err != nil {
		return err
	}

	savePath := path.Join(dirPath, indexName)
	return ioutil.WriteFile(savePath, b[:], 0644)
}
