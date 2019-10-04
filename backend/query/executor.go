package query

import (
	"bogoDB/backend/storage"
	"github.com/pkg/errors"
)

type Executor struct {
	storage *storage.Storage
	catalog *storage.Catalog
}

func NewExecutor(storage *storage.Storage, catalog *storage.Catalog) *Executor{
	return &Executor{
		storage:storage,
		catalog:catalog,
	}
}

func (e *Executor) createTable(q *CreateTableQuery) error {
	return e.catalog.Add(q.Scheme)
}

func (e *Executor) ExecuteMain(q Query) error{
	switch concrete := q.(type) {
		case *CreateTableQuery:
		return e.createTable(concrete)
	}

	return errors.New("failed to execute query")
}