package common

import (
	"os"
	"path/filepath"
)

func GetExecName() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	execName := filepath.Base(execPath)
	// log.Infof("execName:[%v]", execName)
	return execName, nil
}
