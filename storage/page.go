package storage

import (
	"encoding/binary"
	"sync/atomic"
)

const (
	TupleNumber = 32
)

// Page is fixed-sized(4KB) byte chunk as below.

// +----------------+--------+--------+--------------------+
//| PageHeaderData | SlotArray[0] |          |
//+----------------+---+----+---+----+                    |
//|                    |        |                         |
//|                    |  +-----+                         |
//|                    +--+------+                        |
//|                       |      |                        |
//|                       v      v                        |
//|                +----------+-----------------+---------+
//|          <----=+ Item     | Item            | Tuple |
//+----------------+----------+-----------------+---------+

// uint64 is page_id
var currentPgid *uint64
func newPgid() uint64 { return uint64(atomic.AddUint64(currentPgid, 1))}

type Page struct {
	Tuples [TupleNumber]Tuple
}

func NewPage(tuples [TupleNumber]Tuple) *Page{
	return &Page{
		Tuples:tuples,
	}
}

func SerializePage(p *Page)([4096]byte, error){
	var b [4096]byte

	for i, t := range p.Tuples{
		index := i * TupleSize
		tupleBytes, err := SerializeTuple(t)

		if err != nil{
			return b, err
		}

		copy(b[index:index + TupleSize], tupleBytes[:])
	}

	return b, nil
}

func DeserializePage(b [4096]byte) (*Page, error){
	p := &Page{}

	minTxId := binary.BigEndian.Uint64(b[0:8])
	maxTxId := binary.BigEndian.Uint64(b[8:16])

	for i:=0; i<32; i++{
		t := Tuple{
			minTxId:minTxId,
			maxTxId:maxTxId,
			length:b[16],
		}
		copy(t.data[:], b[17:])
		p.Tuples[i] = t
	}

	return p, nil
}
