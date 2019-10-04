package storage

import (
	"bogoDB/backend/meta"
	"github.com/stretchr/testify/assert"
	"log"
	"sync"
	"testing"
)

func TestSaveCatalog(t *testing.T) {
	var s []*meta.Scheme
	ctg := &Catalog{
		Schemes: s,
	}

	ctg.Schemes = append(ctg.Schemes, &meta.Scheme{
		TblName:"users",
		ColTypes:[]meta.ColType{meta.Int},
		ColNames:[]string{"id"},
	})

	if err := SaveCatalog("/tmp/bogodb/", ctg); err != nil{
		log.Fatal(err)
	}
}

func TestLoadCatalog(t *testing.T) {
	var s []*meta.Scheme
	ctg := &Catalog{
		Schemes:s,
	}

	ctg.Schemes = append(ctg.Schemes, &meta.Scheme{
		TblName:"users",
		ColTypes:[]meta.ColType{meta.Int},
		ColNames:[]string{"id"},
	})

	if err := SaveCatalog("/tmp/bogodb/", ctg); err != nil{
		log.Fatal(err)
	}

	out, err := LoadCatalog("/tmp/bogodb/")
	if err != nil{
		log.Fatal(err)
	}

	assert.Equal(t, "users", out.Schemes[0].TblName)
}

func TestAddConcurrenct(t *testing.T){
	var s []*meta.Scheme
	ctg := &Catalog{
		Schemes:s,
	}
	ctg.mutex = &sync.RWMutex{}

	var wg sync.WaitGroup

	for i:=0; i<1000; i++{
		wg.Add(1)
		go func(){
			scheme := &meta.Scheme{TblName:string(i)}
			ctg.Add(scheme)
			wg.Done()
		}()
	}
	wg.Wait()

	assert.Len(t, ctg.Schemes, 1000, "expected %d, but %d", 1000, len(ctg.Schemes))
}