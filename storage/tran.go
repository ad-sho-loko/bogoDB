package storage

import (
	"sync/atomic"
)

type TransactionManager struct {
	// it may overflow someday.
	currentTxid uint64
}

type Transaction struct {
	txid uint64
	tuples []*Tuple
	// modifiedTuples
	storage *Storage
}

func NewTransactionManager() *TransactionManager{
	return &TransactionManager{
		// it should be persisted
		currentTxid:0,
	}
}

func (t *TransactionManager) newTxid() uint64{
	return atomic.AddUint64(&t.currentTxid, 1)
}

func (t *TransactionManager) BeginTransaction(tuples []*Tuple, storage *Storage) *Transaction{
	return &Transaction{
		txid:t.newTxid(),
		tuples:tuples,
		storage:storage,
	}
}

func (t *Transaction) Txid() uint64{
	return t.txid
}

func (t *Transaction) Commit(tableName string){
	for _, tuple := range t.tuples{
		t.storage.InsertTuple(tableName, tuple)
	}
}

func (t *Transaction) Rollback(){
}

/*
type LockManager struct {
	locks map[uint64]sync.RWMutex
}

func (l *LockManager) LockShared(tid uint64){
	l.locks[tid].RLock()
}

func LockShared(tid uint64){
}

func LockExclusive(tid uint64){
}

func (l *LockManager) UnLock(tid uint64){
	l.locks[tid].RUnlock()
}
*/