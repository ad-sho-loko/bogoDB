package query

import (
	"github.com/ad-sho-loko/bogodb/meta"
	"github.com/ad-sho-loko/bogodb/storage"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestExecuteCreateTable(t *testing.T) {
	q := &CreateTableQuery{
		Scheme: &meta.Scheme{
			TblName:  "users",
			ColNames: []string{"id"},
			ColTypes: []meta.ColType{meta.Int},
		},
	}

	ctg := storage.NewEmtpyCatalog()
	e := NewExecutor(nil, ctg, nil)
	assert.False(t, e.catalog.HasScheme("users"))

	if err := e.createTable(q); err != nil {
		log.Fatal(err)
	}

	assert.True(t, e.catalog.HasScheme("users"))
}
