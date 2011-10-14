package locserver

import (
	"sync"
)

type idMaker struct {
	id   int64
	lock sync.Mutex
}

func (i *idMaker) new() int64 {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.id++
	return i.id
}
