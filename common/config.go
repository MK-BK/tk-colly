package common

import (
	"encoding/json"
	"io/ioutil"
)

type Configuration struct {
	DBHost     string `json:"dbHost"`
	DBName     string `json:"dbName"`
	DBPassword string `json:"dbPassword"`
	DBTable    string `json:"dbTable"`
	ListenPort string `json:"listenPort"`
}

func LoadConfig(path string) (*Configuration, error) {
	var configuration Configuration

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &configuration); err != nil {
		return nil, err
	}

	return &configuration, nil
}
