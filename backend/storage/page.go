package storage

import (
	"encoding/binary"
	"sync/atomic"
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
	Pgid uint64
	Type uint8
	// 9byte

	Tuples [31]Tuple
	// 128 * 31 = 3981 byte
	// 128(header) + 3981(data)= 4096
}
// Lower uint32 // pointer of the last line ptr
// Upper uint32 // pointer of the last tuple

func newDataPage() *Page{
	return &Page{
		Pgid:newPgid(),
		Type:1,
	}
}

func newIndexPage() *Page{
	return &Page{
		Pgid:newPgid(),
		Type:2,
	}
}

func SerializePage(p *Page)([4096]byte, error){
	var b [4096]byte

	binary.BigEndian.PutUint64(b[0:8], uint64(p.Pgid))
	b[8] = p.Type

	for i, t := range p.Tuples{
		index := 9 + (i * TupleSize)
		tupleBytes, err := SerializeTuple(t)

		if err != nil{
			return b, err
		}

		copy(b[index:index+TupleSize], tupleBytes[:])
	}

	return b, nil
}

/*
func DeserializePage(b [4096]byte) (*Page, error){
}
*/
