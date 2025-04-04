package common

import (
	"encoding/binary"
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
func Uint64ToBytes(n uint64) []byte {
	bytes := make([]byte, 8) // uint64 占用 8 个字节
	binary.BigEndian.PutUint64(bytes, n)
	return bytes
}

func BytesToUint64(bytes []byte) uint64 {
	return binary.BigEndian.Uint64(bytes)
}
