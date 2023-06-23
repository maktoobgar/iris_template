package app

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	g "service/global"

	"github.com/kataras/iris/v12"
)

const (
	envCloneChildKey    = "IRIS_CLONE_CHILD"
	envCloneChildVal    = "1"
	envCloneChildNumber = "IRIS_CLONE_CHILD_NUM"
)

func IsChild() bool {
	return os.Getenv(envCloneChildKey) == envCloneChildVal
}

func GetChildNumber() string {
	return os.Getenv(envCloneChildNumber)
}

func RunClonesAndServer(app *iris.Application) {
	if IsChild() {
		// use 1 cpu core per child process
		runtime.GOMAXPROCS(1)
		app.Listen(g.CFG.Gateway.IP+":"+g.CFG.Gateway.Port, iris.WithSocketSharding)
		return
	}

	type child struct {
		pid int
		err error
	}

	// create variables
	max := runtime.GOMAXPROCS(g.CFG.ClonesCount)

	if g.CFG.ClonesCount == 0 {
		max = 0
	} else if g.CFG.ClonesCount > max {
		max = g.CFG.ClonesCount
	} else if g.CFG.ClonesCount < max && g.CFG.ClonesCount != -1 {
		max = g.CFG.ClonesCount
	}
	childs := make(map[int]*exec.Cmd)
	channel := make(chan child, max)

	channelShutdownInfo := make(chan any, 1)
	// kill child procs when master exits
	defer func() {
		channelShutdownInfo <- 1
		for _, proc := range childs {
			if err := proc.Process.Kill(); err != nil {
				if !errors.Is(err, os.ErrProcessDone) {
					g.Logger.Error(fmt.Sprintf("clone: failed to kill child: %v\n", err), nil, RunClonesAndServer)
				}
			}
		}
	}()

	// launch child procs
	for i := 0; i < max; i++ {
		cmd := exec.Command(os.Args[0], os.Args[1:]...) //nolint:gosec // It's fine to launch the same process again
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// add fiber clone child flag into child proc env
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("%s=%s", envCloneChildKey, envCloneChildVal),
			fmt.Sprintf("%s=%d", envCloneChildNumber, i+1),
		)
		if err := cmd.Start(); err != nil {
			g.Logger.Error("failed to start a child clone process", nil, RunClonesAndServer)
			return
		}

		// store child process
		pid := cmd.Process.Pid
		childs[pid] = cmd

		// notify master if child crashes
		go func() {
			channel <- child{pid, cmd.Wait()}
		}()
	}

	// Run App
	if max > 0 {
		// return error if child crashes
		select {
		case crashedProcess := <-channel:
			g.Logger.Error(fmt.Sprintf("error: process with %d id crashed", crashedProcess.pid), nil, RunClonesAndServer)
		case <-channelShutdownInfo:
			return
		}
	} else {
		app.Listen(g.CFG.Gateway.IP + ":" + g.CFG.Gateway.Port)
	}
}
