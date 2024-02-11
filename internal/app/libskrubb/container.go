package libskrubb

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	reexec "github.com/docker/docker/pkg/reexec"
)

func init() {
	reexec.Register("containerInit", containerInit)

	if reexec.Init() {
		os.Exit(0)
	}
}

func containerInit() {

	newContainerName := os.Args[1]

	// if err := mountProc(newrootPath); err != nil {
	// 	fmt.Printf("Error mounting /proc - %s\n", err)
	// 	os.Exit(1)
	// }

	claimErr := ClaimContainerSpace(newContainerName)

	if claimErr != nil {
		fmt.Printf("Error running pivot_root - %s\n", claimErr)
		os.Exit(1)
	}

	containerExec()
}

func containerExec() {
	fmt.Println("containerExec called")

	cmd := exec.Command("/bin/sh")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = []string{"PS1=-[skrubb]- # "}

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

// Container is a struct that holds the information about the container
type Container struct {
	Name string
	pid  string
}

// AttatchToContainer creates a new container
func AttatchToContainer(name string) *Container {
	var rootfsPath string
	flag.StringVar(&rootfsPath, "rootfs", "/tmp/skrubb/rootfs", "Path to the root filesystem to use")
	//var netsetgoPath string
	//flag.StringVar(&netsetgoPath, "netsetgo", "/usr/local/bin/netsetgo", "Path to the netsetgo binary")
	flag.Parse()

	cmd := reexec.Command("containerInit", rootfsPath)

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

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting the reexec.Command - %s\n", err)
		os.Exit(1)
	}
	pid := fmt.Sprintf("%d", cmd.Process.Pid)
	fmt.Printf("%s\n", pid)
	// netsetgoCmd := exec.Command(netsetgoPath, "-pid", pid)

	// if err := netsetgoCmd.Run(); err != nil {
	// 	fmt.Printf("Error running netsetgo - %s\n", err)
	// 	os.Exit(1)
	// }

	waitErr := cmd.Wait()
	if waitErr != nil {
		fmt.Printf("%s\n", waitErr)
		os.Exit(1)
	}

	return &Container{
		Name: name,
		pid:  pid,
	}
}
