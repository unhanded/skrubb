package image

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func checkExist(p string) bool {
	_, statErr := os.Stat(p)
	return os.IsExist(statErr)
}

func checkImage(imageFilepath string) bool {
	if !checkExist(imageFilepath) {
		return false
	}
	if !strings.HasSuffix(imageFilepath, "tar.gz") {
		return false
	}
	return true
}

func copyDecompressGzip(filepath string, targetDir string) error {
	command := "tar"
	args := fmt.Sprintf("-C %s -xf %s", targetDir, filepath)
	cmd := exec.Command(command, args)
	startErr := cmd.Start()
	if startErr != nil {
		return startErr
	}
	waitErr := cmd.Wait()
	if waitErr != nil {
		return waitErr
	}
	return nil
}

func CopyImageToContainer(imageFilepath string, containerRoot string) error {
	if !checkImage(imageFilepath) {
		return fmt.Errorf("%s", "File needs to be .tar.gz, and y'know.. exist")
	}
	if !checkExist(containerRoot) {
		return fmt.Errorf("%s", "Container root directory not found")
	}
	return copyDecompressGzip(imageFilepath, containerRoot)
}
