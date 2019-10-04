package storage

import "io/ioutil"

// diskManager is responsible for fetching pages from the disk.
type diskManager struct {
}

func newDiskManager() *diskManager{
	return &diskManager{
	}
}

func (d *diskManager) fetchPage(pid uint64) (*Page, error){
	_, err := ioutil.ReadFile("./tmp/" /* + users/ */+ string(pid))
	if err != nil{
		return nil, err
	}

	/*
	var buf [4096]byte
	copy(buf[:], b)
	page, err := DeserializePage(buf)
	*/

	return nil, nil
}
