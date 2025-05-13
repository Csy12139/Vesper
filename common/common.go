package common

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io"
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

func CalculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate MD5 for %s: %v", filePath, err)
	}
	hashBytes := hash.Sum(nil)
	return fmt.Sprintf("%x", hashBytes), nil
}

func BytesToUint64(bytes []byte) uint64 {
	return binary.BigEndian.Uint64(bytes)
}
