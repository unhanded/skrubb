package pivot

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

func EnterContainer(rootPath string) error {
	sbErr := selfBind(rootPath)
	if sbErr != nil {
		return sbErr
	}

	dumpTarget := filepath.Join(rootPath, ".oldroot/")
	dumpTargetErr := os.MkdirAll(dumpTarget, 0700)

	if dumpTargetErr != nil {
		return dumpTargetErr
	}

	pivotErr := syscall.PivotRoot(rootPath, dumpTarget)
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
