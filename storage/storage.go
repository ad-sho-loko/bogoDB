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
// buffertag -> (table, block_id)
// block_idはpage_idと捉えてよさそう？
func (s *Storage) insertPage(tablename string, pgId uint64, newTuples []*Tuple){
	/*
	isNeedPersist, victim := s.buffer.putPage(pgId, )

	if isNeedPersist {
		// if a victim is dirty, its data must be persisted in the disk now.
		// s.disk.persist(victim)
	}
	*/
}

func (s *Storage) InsertTuple(tablename string, newTuples []*Tuple){
}

func (s *Storage) ReadPage(tableName string, pgid uint64) (*Page, error){
	pg, err := s.buffer.readPage(tableName, 0)

	if err != nil{
		return nil, err
	}

	// if a page exists in the buffer, return that.
	if pg != nil{
		return pg, nil
	}

	pg, err  = s.disk.fetchPage(tableName, pgid)
	if err != nil{
		return nil, err
	}
	s.buffer.putPage(pgid, pg)

	return pg, nil
}
