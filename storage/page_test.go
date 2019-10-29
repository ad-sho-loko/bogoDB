package storage

import (
	"fmt"
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
