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

func (s *Storage) InsertTuple(){
}

func (s *Storage) ReadPage(pgid uint64) (*Page, error){
	pg, err := s.buffer.fetchPage(pgid)

	if err != nil{
		return nil, err
	}

	if pg != nil{
		return pg, nil
	}

	pg, err  = s.disk.fetchPage(pgid)
	if err != nil{
		return nil, err
	}

	return pg, nil
}
