package storage

import "github.com/ad-sho-loko/bogodb/meta"

// Storage is the interface of any manipulation of read/write.
type Storage struct {
	buffer *bufferPool
	disk *diskManager
	home string
}

func NewStorage(home string) *Storage{
	return &Storage{
		buffer:newBufferPool(),
		disk:newDiskManager(),
		home:home,
	}
}

func (s *Storage) insertPage(tableName string){
	pg := NewPage()
	pgid := NewPgid(tableName)
	isNeedPersist, victim := s.buffer.putPage(tableName, pgid, pg)

	if isNeedPersist {
		// if a victim is dirty, its data must be persisted in the disk now.
		if err := s.disk.persist(s.home, tableName, pgid, victim); err != nil{
		}
	}
}

func (s *Storage) InsertTuple(tablename string, t *Tuple){
	for !s.buffer.appendTuple(tablename, t){
		// if not exist in buffer, put a page to lru-cache
		s.insertPage(tablename)
	}
}

func (s *Storage) CreateIndex(indexName string) error{
	btree := meta.NewBTree()
	s.buffer.btree[indexName] = btree
	return nil
}

func (s *Storage) InsertIndex(indexName string, item meta.Item) error{
	btree, err := s.ReadIndex(indexName)
	if err != nil{
		return err
	}

	btree.Insert(item)
	return nil
}

func (s *Storage) ReadIndex(indexName string) (*meta.BTree, error){
	found, index := s.buffer.readIndex(indexName)

	if found{
		return index, nil
	}

	btree, err := s.disk.readIndex(indexName)
	if err != nil{
		return nil, err
	}

	return btree, nil
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

	pg, err = s.disk.fetchPage(s.home, tableName, pgid)

	if err != nil{
		return nil, err
	}
	s.buffer.putPage(tableName, pgid, pg)

	return pg, nil
}

func (s *Storage) Terminate() error{
	items := s.buffer.lru.GetAll()
	for _, item := range items{
		pd := item.(*pageDescriptor)
		if pd.dirty{
			err := s.disk.persist(s.home, pd.tableName, pd.pgid, pd.page)
			if err != nil{
				return err
			}
		}
	}

	for key, val := range s.buffer.btree{
		err := s.disk.writeIndex(s.home, key, val)
		if err != nil{
			return err
		}
	}

	return nil
}