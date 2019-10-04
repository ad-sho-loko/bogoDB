package storage

import (
	"encoding/binary"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"unsafe"
)

func TestPageSize(t *testing.T){
	p := Page{}
	if int(unsafe.Sizeof(p)) > os.Getpagesize(){
		err := fmt.Errorf("bogoDB is compatible with the page whose size should be less than 4KB, but %d", unsafe.Sizeof(p))
		log.Fatal(err)
	}
}

func TestSerialize(t *testing.T){
	p := &Page{}
	p.Pgid = 1
	p.Type = 2

	out, err := SerializePage(p)
	if err != nil{
		log.Fatal(err)
	}

	assert.Equal(t, uint64(1), binary.BigEndian.Uint64(out[0:8]))
	assert.Equal(t, byte(2), out[8])
}

/*
func TestDeserialize(t *testing.T){
	b := [4096]byte{
		// pid
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		// type
		0x01,
		// data
		0xFF,
	}

	p, err := DeserializePage(b)
	if err != nil{
		log.Fatal(err)
	}

	if p.uint64 != 1 || p.Type != 1 || p.Lower != 1 || p.Upper != 1{
		log.Fatal("page header assertion failed")
	}

	if p.Data[0] != 0xFF{
		log.Fatal("page data assertion failed")
	}
}
*/