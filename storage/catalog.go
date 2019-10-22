package storage

import (
	"encoding/json"
	"github.com/ad-sho-loko/bogodb/meta"
	"io/ioutil"
	"path"
	"path/filepath"
	"sync"
)

const catalogName = "catalog.db"

type Catalog struct {
	Schemes []*meta.Scheme
	mutex *sync.RWMutex
}

func NewEmtpyCatalog() *Catalog{
	return &Catalog{
		mutex:&sync.RWMutex{},
	}
}

func LoadCatalog(catalogPath string) (*Catalog, error){
	b, err := ioutil.ReadFile(path.Join(catalogPath, catalogName))
	if err != nil{
		return NewEmtpyCatalog(), nil
	}

	var catalog Catalog
	err = json.Unmarshal(b, &catalog)
	if err != nil{
		return nil, err
	}

	catalog.mutex = &sync.RWMutex{}
	return &catalog, nil
}

// SaveCatalog persists the system catalog as `catalog.db`.
// `catalog.db` has a simple json format like key/value.
func SaveCatalog(dirPath string, c *Catalog) (err error){
	jsonStr, err := json.Marshal(c)
	if err != nil{
		return
	}

 	err = ioutil.WriteFile(filepath.Join(dirPath, catalogName), jsonStr, 0644)
	return
}

// Add is to add the new scheme into a memory.
// Be careful not to persist the disk.
func (c *Catalog) Add(scheme *meta.Scheme) error{
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Schemes = append(c.Schemes, scheme)
	return nil
}

func (c *Catalog) HasScheme(tblName string) bool{
	return c.FetchScheme(tblName) != nil
}

func (c *Catalog) FetchScheme(tblName string) *meta.Scheme{
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// TODO : Fix O(n) to O(1)
	for _, s := range c.Schemes{
		if s.TblName == tblName{
			return s
		}
	}
	return nil
}
