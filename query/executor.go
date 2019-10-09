package query

import (
	"github.com/ad-sho-loko/bogodb/storage"
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

func (s *SeqScan) Scan(store *storage.Storage) []*storage.Tuple{
	var result []*storage.Tuple

	pg, err := store.ReadPage(s.tblName, 0)
	if err != nil{
	}

	for _, t := range pg.Tuples{
		result = append(result, &t)
	}

	return result
}

/*
func (e *Executor) execExpr(expr Expr){
	eq := expr.(*Eq)
}

func (e *Executor) selectTable(q *SelectQuery) error{
	var p Plan

	// from
	tuples := p.scanners.Scan(e.storage)

	// where
	q.Where

	// select
	//q.Col.n

	return nil
}
*/

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