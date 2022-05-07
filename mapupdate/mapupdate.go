package mapupdate

import (
	"sync"
)

type Comparable[T any] interface {
	NE(obj T) bool
}

func Update[R, T Comparable[R]](origMap, newMap map[string]R, mapLock *sync.RWMutex) (set map[string]R, del map[string]R) {
	set = make(map[string]R)
	del = make(map[string]R)

	for k, v := range newMap {
		ex, found := origMap[k]
		if !found || ex.NE(v) {
			set[k] = v
		}
	}

	for k, v := range origMap {
		if _, found := newMap[k]; !found {
			del[k] = v
		}
	}

	mapLock.Lock()
	defer mapLock.Unlock()

	for k, v := range set {
		origMap[k] = v
	}

	for k := range del {
		delete(origMap, k)
	}

	return
}
