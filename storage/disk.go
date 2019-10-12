package storage

import (
	"io/ioutil"
	"os"
	"path"
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

func (d *diskManager) fetchPage(tableName string, pid uint64) (*Page, error){
	pagePath := path.Join(tableName, string(pid))

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

func (d *diskManager) persist(tableName string, page *Page) error{
	b, err := SerializePage(page)

	if err != nil{
		return err
	}

	return ioutil.WriteFile(tableName, b[:], 0644)
}