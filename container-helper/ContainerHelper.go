package main

import (
 "fmt"
 ps "github.com/mitchellh/go-ps"
 "os"
 "strconv"
)

// NewContainer instantiates a default password store
func NewContainerUtil() PidCache {
	return PidCache{
		pidCache:= NewPidCache()
	}
}

type ContainerUtil struct {
	pidCache PidCache
}

func (cu ContainerUtil) getContainerId(pid string) (string, error) {

 cid, err := cu.pidCache.Get(pid)

 if (cid != nil && err == nil ) {
   return cid
 }


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
 cu.pidCache.Set(pid, cid)

return cid

}
