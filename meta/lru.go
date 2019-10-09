package meta

import (
	"container/list"
	"sync"
)

// Lru is the thread-safe LRU cache.
type Lru struct {
	cap int
	evictList *list.List
	items map[interface{}]*list.Element
	mutex sync.RWMutex
}

type entry struct {
	key interface{}
	value interface{}
}

func NewLru(cap int) *Lru {
	return &Lru{
		cap:cap,
		evictList:list.New(),
		items:make(map[interface{}]*list.Element),
	}
}

func (l *Lru) Insert(key, value interface{}) interface{} {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	var victim interface{}
	ent := &entry{key, value}
	elm := l.evictList.PushFront(ent)
	l.items[key] = elm

	if l.needEvict(){
		victim = l.evictList.Back()
		l.removeOldest()
	}

	return victim
}

func (l *Lru) Get(key interface{}) interface{}{
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	if elm, ok := l.items[key]; ok{
		l.evictList.MoveToFront(elm)
		return elm.Value.(*entry).value
	}

	return nil
}

func (l *Lru) Len() int{
	return l.evictList.Len()
}

func (l *Lru) removeOldest(){
	elm := l.evictList.Back()
	if elm != nil{
		l.evictList.Remove(elm)
		delete(l.items, elm.Value.(*entry).key)
	}
}

func (l *Lru) needEvict() bool{
	return l.Len() > l.cap
}
