package storage

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

// diskManager is responsible for fetching pages from the disk.
type diskManager struct {
}

func newDiskManager() *diskManager{
	return &diskManager{
	}
}

func (d *diskManager) toPid(tid uint64) uint64{
	return tid / TupleNumber
}

func (d *diskManager) fetchPageByTid(tableName string, tid uint64) (*Page, error){
	pid := d.toPid(tid)
	return d.fetchPage(tableName, pid)
}

func (d *diskManager) fetchPage(tableName string, pgid uint64) (*Page, error){
	fileName := strconv.FormatUint(pgid, 10)
	pagePath := path.Join(tableName, fileName)

	_, err := os.Stat(pagePath)
	if os.IsNotExist(err){
		return nil, err
	}

	bytes, err := ioutil.ReadFile(pagePath)
	if err != nil{
		return nil, err
	}

	var b [4096]byte
	copy(b[:], bytes)
	return DeserializePage(b)
}

func (d *diskManager) persist(tableName string, pgid uint64, page *Page) error{
	b, err := SerializePage(page)

	if err != nil{
		return err
	}

	fileName := strconv.FormatUint(pgid, 10)
	savePath := path.Join(tableName, fileName)
	return ioutil.WriteFile(savePath, b[:], 0644)
}