package storage

import (
	"io/ioutil"
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
	return tid
}

func (d *diskManager) fetchPage(tableName string, tid uint64) (*Page, error){
	pid := d.toPid(tid)

	bytes, err := ioutil.ReadFile(path.Join(tableName, string(pid)))
	if err != nil{
		return nil, err
	}

	var b [4096]byte
	copy(b[:], bytes)
	return DeserializePage(b)
}
