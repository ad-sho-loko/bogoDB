package meta

import "sync/atomic"

type Transaction struct {
	txid uint64
	// modifiedTuples
}

// it may overflow someday.
var currentTxid uint64 = 0

func newTxid() uint64{
	return atomic.AddUint64(&currentTxid, 1)
}

func BeginTransaction() *Transaction{
	return &Transaction{
		txid:newTxid(),
	}
}

/*
func Commit(tran *Transaction){
}

func Rollback(tran *Transaction){
}
*/
