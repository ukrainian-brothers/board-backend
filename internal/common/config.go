package common

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
}

type Config struct {
	Postgres PostgresConfig `json:"postgres_config"`
}

func NewConfigFromFile(fileName string) (*Config, error) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	cfg := Config{}
	err = json.Unmarshal(byteValue, &cfg)
	return &cfg, err
}
