package storage

// Storage is the interface of any manipulation of read/write.
type Storage struct {
	buffer *bufferPool
	disk *diskManager
}

func NewStorage() *Storage{
	return &Storage{
		buffer:newBufferPool(),
		disk:newDiskManager(),
	}
}

// postgresqlはbuffertag単位で保存する.
func (s *Storage) insertPage(tablename string, pgid uint64, newTuples *[32]Tuple){
	pg := &Page{
		Tuples:*newTuples,
	}

	isNeedPersist, victim := s.buffer.putPage(pgid, pg)

	if isNeedPersist {
		// if a victim is dirty, its data must be persisted in the disk now.

		err := s.disk.persist(tablename, victim)
		if err != nil{

		}
	}
}

func (s *Storage) InsertTuple(tablename string, newTuples *Tuple){
}

func (s *Storage) ReadTuple(tableName string, tid uint64) (*Tuple, error){
	pgid := s.buffer.toPid(tid)

	pg, err := s.readPage(tableName, pgid)
	if err != nil{
		return nil, err
	}

	return &pg.Tuples[tid % TupleNumber], nil
}

func (s *Storage) readPage(tableName string, pgid uint64) (*Page, error){
	pg, err := s.buffer.readPage(tableName, pgid)

	if err != nil{
		return nil, err
	}

	// if a page exists in the buffer, return that.
	if pg != nil{
		return pg, nil
	}

	pg, err = s.disk.fetchPage(tableName, pgid)

	if err != nil{
		return nil, err
	}
	s.buffer.putPage(pgid, pg)

	return pg, nil
}
