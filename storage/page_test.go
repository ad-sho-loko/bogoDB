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

func TestPageSize(t *testing.T) {
	p := Page{}
	if int(unsafe.Sizeof(p)) > os.Getpagesize() {
		err := fmt.Errorf("bogoDB is compatible with the page whose size should be less than 4KB, but %d", unsafe.Sizeof(p))
		log.Fatal(err)
	}
}

func TestSerialize(t *testing.T) {
	p := &Page{}

	out, err := SerializePage(p)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, uint64(1), binary.BigEndian.Uint64(out[0:8]))
	assert.Equal(t, byte(2), out[8])
}
