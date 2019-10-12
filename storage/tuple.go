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
	length uint8
	data [111]byte
}

func NewTuple(minTxId uint64, values []interface{}) *Tuple{
	var b [111]byte

	i := 0
	length := uint8(0)
	for _, v := range values{
		switch concrete := v.(type) {
		case int:
			binary.BigEndian.PutUint32(b[i:i+4], uint32(concrete))
			length++
			i+=4
		case string:
			// b[i] = uint8(len(concrete))
			// utf32 := []byte(concrete)
			// utf2
			length++
		}
	}

	return &Tuple{
		minTxId:minTxId,
		maxTxId:minTxId,
		length:length,
		data:b,
	}
}

func SerializeTuple(t Tuple) ([TupleSize]byte, error){
	var b [TupleSize]byte

	binary.BigEndian.PutUint64(b[0:8], t.minTxId)
	binary.BigEndian.PutUint64(b[8:16], t.maxTxId)
	b[16] = t.length
	copy(b[17:], t.data[:])

	return b, nil
}

func DeserializeTuple(b [TupleSize]byte) (Tuple, error){
	var t Tuple

	t.minTxId = binary.BigEndian.Uint64(b[0:8])
	t.maxTxId = binary.BigEndian.Uint64(b[8:16])
	t.length = b[16]
	copy(t.data[:], b[17:])

	return t, nil
}

func (t *Tuple) IsUnused() bool{
	// If minTxId is zero, it's an empty tuple.
	return t.minTxId == 0
}