package query

import (
	"fmt"
	"github.com/ad-sho-loko/bogodb/storage"
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

// SeqScan
func (s *SeqScan) Scan(store *storage.Storage) []*storage.Tuple{
	var result []*storage.Tuple

	for i:=uint64(0);; i++{
		t, err := store.ReadTuple(s.tblName, i)
		if err != nil{
			// if no more pages, finish reading tuples.
			break
		}

		if t.IsUnused(){
			// if no more tuples in page, finish reading tuples.
			break
		}

		result = append(result, t)
	}
	return result
}

// IndexScan
func (s *IndexScan) Scan(store *storage.Storage) []*storage.Tuple{
	var result []*storage.Tuple
	// idx := s.tblName + s.indexCol
	// btree, _ := store.ReadIndex(idx)
	// items := btree.Get()
	return result
}

func (e *Executor) where(tuples []*storage.Tuple, tableName string, where []Expr) []*storage.Tuple{
	var filtered []*storage.Tuple
	for _, w := range where{
		left := w.(*Eq).left.(*Lit)
		right := w.(*Eq).right.(*Lit)

		for _, t := range tuples{
			if t.GetStr(tableName, left.v) == right.v{
				// filtered = append(filtered, t)
			}
		}
	}

	return filtered
}

func (e *Executor) selectTable(q *SelectQuery, p *Plan, tran *storage.Transaction) error{
	// from, where
	tuples := p.scanners.Scan(e.storage)
	filteredTuple := e.where(tuples, q.From[0].Name, q.Where)

	// consider transactions.
	for _, t := range filteredTuple{
		if tran == nil || t.CanSee(tran){
			fmt.Println(t)
		}
	}

	// select

	return nil
}

func (e *Executor) insertTable(w *InsertQuery, tran *storage.Transaction) error{
	inTransaction := tran != nil

	if !inTransaction {
		tran = e.beginTransaction()
	}

	t := storage.NewTuple(tran.Txid(), w.Values)
	e.storage.InsertTuple(w.Table.Name, t)
	e.storage.InsertIndex(w.Index, t)

	if !inTransaction{
		e.commitTransaction(tran)
	}

	return nil
}

func (e *Executor) updateTable(q *UpdateQuery, p *Plan, tran *storage.Transaction) error{
	// set, where
	// tuples := p.scanners.Scan(e.storage)
	// filteredTuple := e.where(tuples, q.Where)

	return nil
}

func (e *Executor) createTable(q *CreateTableQuery) error {
	err := e.catalog.Add(q.Scheme)
	if err != nil{
		return err
	}

	return e.storage.CreateIndex(q.Scheme.TblName + "_" + q.Scheme.PrimaryKey)
}

func (e *Executor) beginTransaction() *storage.Transaction{
	return e.tranManager.BeginTransaction()
}

func (e *Executor) commitTransaction(tran *storage.Transaction){
	e.tranManager.Commit(tran)
}

func (e *Executor) abortTransaction(tran *storage.Transaction){
	e.tranManager.Abort(tran)
}

func (e *Executor) ExecuteMain(q Query, p *Plan, tran *storage.Transaction) error{
	switch concrete := q.(type) {
	case *BeginQuery:
		e.beginTransaction()
		return nil
	case *CommitQuery:
		e.commitTransaction(tran)
		return nil
	case *AbortQuery:
		e.abortTransaction(tran)
		return nil
	case *CreateTableQuery:
		return e.createTable(concrete)
	case *InsertQuery:
		return e.insertTable(concrete, tran)
	case *UpdateQuery:
		return e.updateTable(concrete, p, tran)
	case *SelectQuery:
		return e.selectTable(concrete, p, tran)
	}

	return errors.New("failed to execute query")
}