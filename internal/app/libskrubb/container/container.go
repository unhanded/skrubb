package container

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	reexec "github.com/docker/docker/pkg/reexec"
	"github.com/unhanded/skrubb/internal/app/libskrubb/pivot"
)

func init() {
	reexec.Register("containerInit", containerInit)

	if reexec.Init() {
		os.Exit(0)
	}
}

func mountNewProc(rootPath string) error {
	literallyProc := "proc"
	targetPath := filepath.Join(rootPath, "/proc")

	os.MkdirAll(targetPath, 0755)
	mountErr := syscall.Mount(literallyProc, targetPath, literallyProc, uintptr(0), "")
	if mountErr != nil {
		return mountErr
	}

	return nil
}

func containerInit() {

	containersDirPath := os.Args[1]
	containerName := os.Args[2]
	containerRootPath := filepath.Join(containersDirPath, containerName+"/")
	claimErr := pivot.EnterContainer(containerRootPath)
	if claimErr != nil {
		fmt.Printf("Error claiming space for container - %s\n", claimErr)
		os.Exit(1)
	}

	err := mountNewProc(containerRootPath)
	if err != nil {
		fmt.Printf("Error mounting /proc - %s\n", err)
		os.Exit(1)
	}
	containerExec()
}

func containerExec() {
	containerName := os.Args[2]
	cmd := exec.Command("/bin/sh")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	ps1 := fmt.Sprintf("PS1=%s]- # ", containerName)
	cmd.Env = []string{ps1}

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error %s in namespaceExec\n", err)
		os.Exit(1)
	}
}

type SysProcIDMap struct {
	ContainerID int // Container ID.
	HostID      int // Host ID.
	Size        int // Size.
}

func MakeContainer() error {
	var containersDirPath string
	flag.StringVar(&containersDirPath, "containersDirPath", "/tmp/skrubb/containers/", "Path to the directory where all the containers go")
	var containerName string
	flag.StringVar(&containerName, "containerName", "skrubbox", "Name of the container")

	flag.Parse()

	cmd := reexec.Command("containerInit", containersDirPath, containerName)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

	startErr := cmd.Start()
	if startErr != nil {
		fmt.Printf("Error starting the reexec.Command - %s\n", startErr)
		return startErr
	}
	pid := fmt.Sprintf("PID: %d", cmd.Process.Pid)
	fmt.Printf("%s\n", pid)

	waitErr := cmd.Wait()
	if waitErr != nil {
		return waitErr
	}

	return nil
}
