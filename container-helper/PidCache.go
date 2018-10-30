package ContainerHelper

import (
	"errors"
	"sync"
	"time"
	ps "github.com/mitchellh/go-ps"
)

// NewPidCache instantiates a default password store
func NewPidCache() PidCache {
	return PidCache{
		lock:     &sync.RWMutex{},
		pidCache: make(map[int]string),
	}
}

type PidCache struct {
	lock     *sync.RWMutex
	pidCache map[int]string
}

//cleanupLoop that keeps calling cleanupCache() in a loop every scheduled interval
func (pc PidCache) cleanupLoop() {
	//Hardcoded to run once a minute. Will make this configurable
	for range time.Tick(time.Minute * 1) {
		pc.cleanupCache()
	}

}

// This function is meant to clean up the pids that have entries in the cache by the
// corresponding processes have exited
func (pc PidCache) cleanupCache() {
	for pid := range pc.pidCache {
			 p, err := ps.FindProcess(pid)

		 	if err != nil || p == nil {
				pc.Delete(pid)
		 	}
	 }
}

// Set stores the PID and the container Id in the map. If the container id is not available, it will be stored as "non-container".
// If the process tree is killed by the time the container ID is fetched, it will be marked as "killed" as a hint to be purged.
func (pc PidCache) Set(pid int, cid string) error {
	pc.lock.Lock()
	defer pc.lock.Unlock()
	pc.pidCache[pid] = cid
	return nil
}

// Get retrieves the cid, given the pid.
func (pc PidCache) Get(pid int) (string, error) {
	pc.lock.RLock()
	cid, ok := pc.pidCache[pid]
	pc.lock.RUnlock()

	if ok {
		return cid, nil
	} else {
		return "", errors.New("PID not found in cache")
	}
}

func (pc PidCache) Delete(pid int) error {
	pc.lock.Lock()
	defer pc.lock.Unlock()
	delete(pc.pidCache, pid)
	return nil
}

func (pc PidCache) Init() error {
	//Launch a seperate cleanup job thread
	go pc.cleanupLoop()
	return nil
}
