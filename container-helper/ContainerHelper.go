package ContainerHelper

import (
	"fmt"
	ps "github.com/mitchellh/go-ps"
	"io/ioutil"
	"regexp"
	"strings"
	"os"
	"log"
)

var l = log.New(os.Stdout, "", 0)

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
		l.Println(fmt.Sprintf("Error, process not found: %d", pid))
		fmt.Println("Error, process not found")
		return "", err
	}

	not_init := true
	var ncid = 0 //numeric container process id
	var scid, uuid string = "", "" //string process ID with container uuid
	for not_init {

		if strings.HasPrefix(p.Executable(), "docker-containe") {
			l.Println(fmt.Sprintf("Found container shim for %d", pid))
			ncid = p.Pid()
			//find the container uuid from the pid, the container uuid is associated wit the root process not the container shim process
			b, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cgroup", pid))
			if err != nil || b == nil {
				l.Println("Error, container uuid not found")
				//at least record the container process id
				scid = fmt.Sprintf("%d:", ncid)
			} else {
				//parse uuid from file contents
				re := regexp.MustCompile("^\\d+:\\w+:/docker/([0-9a-f]{64})")
				rs := re.FindStringSubmatch(string(b))
				if rs != nil {
					uuid = rs[1]
					scid =  uuid
					l.Println(fmt.Sprintf("UUID %s for %d", scid, pid))
				} else {
					l.Println(fmt.Sprintf("Error: file /proc/%d/cgroup did not match regex. (%s)", pid, string(b)))
				}
			}
			cu.pidCache.Set(pid, scid)
			return scid, nil
		}

		p, err = ps.FindProcess(p.PPid())

		if p == nil || err != nil {
			l.Println(fmt.Sprintf("No parent process for %d", pid))
			fmt.Println("Error : ", err)
			return "-1", nil
		}

		if 1 == p.Pid() {
			not_init = false
		}
	}

	l.Println(fmt.Sprintf("Container shim not found for process for %d", pid))
	cu.pidCache.Set(pid, "0")
	return "0", nil
}
