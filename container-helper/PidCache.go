// NewPidCache instantiates a default password store
func NewPidCache() PidCache {
	return PidCache{
		lock:        &sync.RWMutex{},
		pidCache: make(map[string]string),
	}
}

type PidCache struct {
	lock        *sync.RWMutex
	pidCache map[string]string
}

// Set stores the PID and the container Id in the map. If the container id is not available, it will be stored as "non-container".
// If the process tree is killed by the time the container ID is fetched, it will be marked as "killed" as a hint to be purged.
func (pc PidCache) Set(pid string, cid string) error {
		pc.lock.Lock()
		defer pc.lock.Unlock()
		pc.pidCache[pid] = cid
}

// Get retrieves the cid, given the pid.
func (pc PidCache) Get(pid string) (string, error) {
	pc.lock.RLock()
	cid, ok := pc.pidCache[pid]
	pc.lock.RUnlock()
	if ok {
		return cid, nil
	}
}
