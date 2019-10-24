package storage

const (
	TupleNumber = 32
)

// Page is fixed-sized(4KB) byte.
type Page struct {
	Tuples [TupleNumber]Tuple
}

func NewPage() *Page {
	return &Page{
		Tuples: [TupleNumber]Tuple{},
	}
}

func NewPgid(tableName string) uint64 {
	// FIXME
	return 0
}

func SerializePage(p *Page) ([4096]byte, error) {
	var b [4096]byte

	for i, t := range p.Tuples {
		tupleBytes, err := SerializeTuple(&t)

		if err != nil {
			return b, err
		}

		copy(b[i*128:i*128+128], tupleBytes[:])
	}

	return b, nil
}

func DeserializePage(b [4096]byte) (*Page, error) {
	p := &Page{}

	for i := 0; i < 32; i++ {
		var in [128]byte
		copy(in[:], b[i*128:i*128+128])
		t, err := DeserializeTuple(in)

		if err != nil {
			return nil, err
		}

		p.Tuples[i] = *t
	}

	return p, nil
}
