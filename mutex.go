package redlock

import (
	"sync"

	"github.com/go-redsync/redsync/v4"
)

// Mutexes contains slice of mutex
type Mutexes []*redsync.Mutex

// Lock locks mutexes. In case it returns an error on failure, you may retry to acquire the lock by calling this method again
func (m *Mutexes) Lock() (err error) {
	for _, v := range *m {
		if err = v.Lock(); err != nil {
			return
		}
	}

	return
}

// Unlock unlocks mutexes and return its error.
func (m *Mutexes) Unlock() (err error) {
	var wg sync.WaitGroup
	for _, v := range *m {
		wg.Add(1)
		go func(v *redsync.Mutex) {
			_, err = v.Unlock()
			wg.Done()
		}(v)
	}
	wg.Wait()

	return
}
