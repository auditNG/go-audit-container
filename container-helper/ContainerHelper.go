package main

import (
 "fmt"
 ps "github.com/mitchellh/go-ps"
 "os"
 "strconv"
)

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

func getContainerId(pid string) (string, error) {

 pid, err := strconv.Atoi(os.Args[1])

 if err != nil {
   fmt.Println("Bad process ID supplied")
   return err
 }

 // at this stage the Processes related functions found in Golang's OS package
 // is no longer sufficient, we will use Mitchell Hashimoto's https://github.com/mitchellh/go-ps
 // package to find the application/executable/binary name behind the process ID.

 p, err := ps.FindProcess(pid)

 if err != nil {
   fmt.Println("Error : ", err)
   os.Exit(-1)
 }

 fmt.Println("Process ID : ", p.Pid())
 fmt.Println("Parent Process ID : ", p.PPid())
 fmt.Println("Process ID binary name : ", p.Executable())

 not_init := true
 container := false
 for not_init {

   if (p.Executable() == "docker-containe") {
     container = true
     fmt.Println("The container ID is: ", p.Pid())
     break
   }

   p, err = ps.FindProcess(p.PPid())

   if err != nil {
     fmt.Println("Error : ", err)
     os.Exit(-1)
   }

   if(1 == p.Pid()) {
     not_init = false
   }
}

if ( container == true ) {
 fmt.Println("This process runs in a container")
} else {
 fmt.Println("This process does not run in a container")
}

}
