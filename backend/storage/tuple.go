package storage

import (
	"encoding/binary"
)

const TupleSize = 128

// Tuple is the actual user's data stored in a `Page`.
// A general tuple should be variable, but bodoDB's one is fixed in 128 byte
type Tuple struct {
	minTxId uint64 // txId when inserted
	maxTxId uint64 // txId when updated
	data [112]byte
}

func SerializeTuple(t Tuple) ([TupleSize]byte, error){
	var b [TupleSize]byte

	binary.BigEndian.PutUint64(b[0:8], t.minTxId)
	binary.BigEndian.PutUint64(b[8:16], t.maxTxId)
	copy(b[16:], t.data[:])

	return b, nil
}

func DeserializeTuple(b [TupleSize]byte) (Tuple, error){
	var t Tuple

	t.minTxId = binary.BigEndian.Uint64(b[0:8])
	t.maxTxId = binary.BigEndian.Uint64(b[8:16])
	copy(t.data[:], b[16:])

	return t, nil
}
