package query

import (
	"fmt"
	"github.com/ad-sho-loko/bogodb/storage"
	"github.com/go-ffmt/ffmt"
	"github.com/pkg/errors"
)

type Executor struct {
	storage *storage.Storage
	catalog *storage.Catalog
	tranManager *storage.TransactionManager
}

func NewExecutor(storage *storage.Storage, catalog *storage.Catalog, tranManger *storage.TransactionManager) *Executor{
	return &Executor{
		storage:storage,
		catalog:catalog,
		tranManager:tranManger,
	}
}

func (s *SeqScan) Scan(store *storage.Storage) []*storage.Tuple{
	var result []*storage.Tuple

	for i:=uint64(0);; i++{
		t, err := store.ReadTuple(s.tblName, i)
		if err != nil{
			// if no more tuples, read tuples end.
			break
		}
		result = append(result, t)
	}

	return result
}

/*
func (e *Executor) execExpr(expr Expr){
	eq := expr.(*Eq)
}
*/

func (e *Executor) selectTable(q *SelectQuery) error{
	var p Plan

	// from
	tuples := p.scanners.Scan(e.storage)

	// where
	// q.Where

	// select
	// q.Col.n

	// print!
	ffmt.Puts(tuples)

	return nil
}

func (e *Executor) insertTable(w *InsertQuery) error{
	if !e.catalog.HasScheme(w.Table.Name){
		return fmt.Errorf("insert failed : `%s` doesn't exists", w.Table.Name)
	}

	var tuples []*storage.Tuple
	tx := e.tranManager.BeginTransaction(tuples, e.storage)
	// t := storage.NewTuple(tx.Txid(), w.Values)
	tx.Commit(w.Table.Name)
	return nil

}

func (e *Executor) createTable(q *CreateTableQuery) error {
	return e.catalog.Add(q.Scheme)
}

func (e *Executor) ExecuteMain(q Query) error{
	switch concrete := q.(type) {
	case *CreateTableQuery:
		return e.createTable(concrete)
	case *InsertQuery:
		return e.insertTable(concrete)
	}

	return errors.New("failed to execute query")
}