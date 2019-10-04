package storage

import (
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestSerializeTuple(t *testing.T) {
	in := [112]byte{3,4,5}

	tuple := Tuple{
		minTxId:1,
		maxTxId:2,
		data:in,
	}

	out, err := SerializeTuple(tuple)
	if err != nil{
		log.Fatal(err)
	}

	assert.Equal(t, uint64(1), binary.BigEndian.Uint64(out[0:8]))
	assert.Equal(t, uint64(2), binary.BigEndian.Uint64(out[8:16]))
	assert.Equal(t, byte(3), out[16])
	assert.Equal(t, byte(4), out[17])
	assert.Equal(t, byte(5), out[18])
}

func TestDeserializeTuple(t *testing.T) {
	in := [128]byte{
		// minTxId
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		// maxTxId
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
		// data
		0x03, 0x04, 0x05,
	}

	out, err := DeserializeTuple(in)
	if err != nil{
		log.Fatal(err)
	}

	assert.Equal(t, uint64(1), out.minTxId)
	assert.Equal(t, uint64(2), out.maxTxId)
	assert.Equal(t, byte(3), out.data[0])
	assert.Equal(t, byte(4), out.data[1])
	assert.Equal(t, byte(5), out.data[2])
}