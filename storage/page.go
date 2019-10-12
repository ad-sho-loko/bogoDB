package storage

import (
	"encoding/binary"
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

type Page struct {
	Tuples [TupleNumber]Tuple
}

func NewPage() *Page{
	return &Page{
		Tuples:[TupleNumber]Tuple{},
	}
}

func NewPgid(tableName string) uint64{
	return 0
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
