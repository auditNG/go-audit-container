package ContainerHelper

import (
	"errors"
	"sync"
)

// NewPidCache instantiates a default password store
func NewPidCache() PidCache {
	return PidCache{
		lock:     &sync.RWMutex{},
		pidCache: make(map[int]int),
	}
}

type PidCache struct {
	lock     *sync.RWMutex
	pidCache map[int]int
}

// Set stores the PID and the container Id in the map. If the container id is not available, it will be stored as "non-container".
// If the process tree is killed by the time the container ID is fetched, it will be marked as "killed" as a hint to be purged.
func (pc PidCache) Set(pid int, cid int) error {
	pc.lock.Lock()
	defer pc.lock.Unlock()
	pc.pidCache[pid] = cid
	return nil
}

// Get retrieves the cid, given the pid.
func (pc PidCache) Get(pid int) (int, error) {
	pc.lock.RLock()
	cid, ok := pc.pidCache[pid]
	pc.lock.RUnlock()

	if ok {
		return cid, nil
	} else {
		return -1, errors.New("PID not found in cache")
	}
}
