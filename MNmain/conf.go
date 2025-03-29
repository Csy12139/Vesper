package MNmain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type LogConfig struct {
	LogDir        string `json:"logDir"`
	MaxFileSizeMb int    `json:"maxFileSizeMb"`
	MaxFileNum    int    `json:"maxFileNum"`
	LogLevel      string `json:"logLevel"`
}

type Config struct {
	MNAddr    string    `json:"mn_addr"`
	MNNetwork string    `json:"mn_network"`
	Log       LogConfig `json:"log"`
}

var GlobalConfig Config

func loadConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open config file: %w", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("could not read config file: %w", err)
	}

	err = json.Unmarshal(bytes, &GlobalConfig)
	if err != nil {
		return fmt.Errorf("could not parse config file: %w", err)
	}
	return nil
}
