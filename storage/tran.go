package storage

import (
	"sync/atomic"
)

type TransactionManager struct {
	clogs       map[uint64]*Transaction // should be thread-safe
	currentTxid uint64                  // it may overflow someday.
}

type Transaction struct {
	txid  uint64
	state TransactionState
}

type TransactionState int

const (
	Commited TransactionState = iota + 1
	InProgress
	Abort
)

func NewTransactionManager() *TransactionManager {
	return &TransactionManager{
		clogs: make(map[uint64]*Transaction),
		// FIXME : it should be persisted when a server shutdown.
		currentTxid: 0,
	}
}

func (t *TransactionManager) newTxid() uint64 {
	// txid starts from 1.
	return atomic.AddUint64(&t.currentTxid, 1)
}

func (t *TransactionManager) BeginTransaction() *Transaction {
	txid := t.newTxid()
	tx := &Transaction{
		txid:  txid,
		state: InProgress,
	}
	t.clogs[txid] = tx
	return tx
}

func (t *TransactionManager) Commit(tran *Transaction) {
	tran.state = Commited
}

func (t *TransactionManager) Abort(tran *Transaction) {
	tran.state = Abort
}

func (t *Transaction) Txid() uint64 {
	return t.txid
}