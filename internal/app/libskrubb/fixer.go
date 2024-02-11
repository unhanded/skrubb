package libskrubb

import (
	"os"
	"path/filepath"
	"syscall"
)

func selfBind(rootpath string) error {
	err := syscall.Mount(
		rootpath,
		rootpath,
		"",
		syscall.MS_BIND|syscall.MS_REC,
		"")
	if err != nil {
		return err
	}
	return nil
}

func doPivot(newRoot string, oldDump string) error {
	err := syscall.PivotRoot(newRoot, oldDump)
	if err != nil {
		return err
	}
	return nil
}

func doChdir(p string) error {
	err := os.Chdir(p)
	if err != nil {
		return err
	}
	return nil
}

func dirRemoval(p string) error {
	umountErr := syscall.Unmount(p, syscall.MNT_DETACH)
	if umountErr != nil {
		return umountErr
	}
	removeErr := os.RemoveAll(p)
	if removeErr != nil {
		return removeErr
	}
	return nil
}

func makeContainerPath(containerName string) string {
	return filepath.Join("/tmp/skrubb/rootfs")
}

func ClaimContainerSpace(rootPath string) error {
	selfBind(rootPath)

	dumpTarget := filepath.Join(rootPath, "/.oldroot/")
	dumpTargetErr := os.MkdirAll(dumpTarget, 0700)
	if dumpTargetErr != nil {
		return dumpTargetErr
	}
	pivotErr := doPivot(rootPath, dumpTarget)
	if pivotErr != nil {
		return pivotErr
	}
	chdirErr := os.Chdir("/")
	if chdirErr != nil {
		return chdirErr
	}
	drErr := dirRemoval("/.oldroot/")
	return drErr
}
