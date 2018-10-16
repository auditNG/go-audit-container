package ContainerHelper

import (
	"fmt"
	ps "github.com/mitchellh/go-ps"
	"io/ioutil"
	"regexp"
)

type ContainerUtil struct {
	pidCache PidCache
}

// NewContainer instantiates a default password store
func NewContainerUtil() ContainerUtil {
	return ContainerUtil{
		NewPidCache(),
	}
}

func (cu ContainerUtil) Init() error {
	return cu.pidCache.Init()
}

func (cu ContainerUtil) GetContainerId(pid int) (string, error) {

	cid, err := cu.pidCache.Get(pid)

	if err == nil {
		return cid, nil
	}

	p, err := ps.FindProcess(pid)

	if err != nil || p == nil {
		fmt.Println("Error, process not found")
		return "", err
	}

	not_init := true
	var ncid = 0 //numeric container process id
	var scid, uuid string = "", "" //string process ID with container uuid
	for not_init {

		if p.Executable() == "docker-containe" {
			ncid = p.Pid()
			//find the container uuid from the pid, the container uuid is associated wit the root process not the container shim process
			b, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cgroup", pid))
			if err != nil || b == nil {
				fmt.Println("Error, container uuid not found")
				//at least record the prociss id
				scid = fmt.Sprintf("%d", ncid)
			} else {
				//parse uuid from file contents
				re := regexp.MustCompile("^\\d+:memory:/docker/([0-9a-f]{64})")
				rs := re.FindStringSubmatch(string(b))
				if rs != nil {
					uuid = rs[1]
					scid = uuid
				}
			}
			cu.pidCache.Set(pid, scid)
			return scid, nil
		}

		p, err = ps.FindProcess(p.PPid())

		if p == nil || err != nil {
			fmt.Println("Error : ", err)
        		return "", nil

		}

		if 1 == p.Pid() {
			not_init = false
		}
	}

	cu.pidCache.Set(pid, "0:")
	return "0:", nil
}
