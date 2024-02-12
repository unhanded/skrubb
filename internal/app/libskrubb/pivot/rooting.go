package pivot

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
)

type SkrubbRoot struct {
	rootPath string
}

func newSkrubbRootPath(containerName string) (string, error) {
	parent := "/tmp/skrubb/"
	var randNum = make([]byte, 5)
	_, rErr := rand.Reader.Read(randNum)
	if rErr != nil {
		return "", nil
	}
	rHash := hex.EncodeToString(randNum)
	cName := rHash + "__" + containerName + "/"
	newRoot := filepath.Join(parent, cName)
	return newRoot, nil
}

func mkdirAtPath(fp string) error {
	_, statErr := os.Stat(fp)
	if os.IsExist(statErr) {
		return statErr
	}
	mkdirErr := os.MkdirAll(fp, 0700)
	if mkdirErr != nil {
		return mkdirErr
	}
	return nil
}

func NewSkrubbRoot(containerName string) (string, error) {
	newPath, err := newSkrubbRootPath(containerName)
	if err != nil {
		return "", err
	}

	mkdirErr := mkdirAtPath(newPath)
	if mkdirErr != nil {
		return "", mkdirErr
	}
	return newPath, nil
}
