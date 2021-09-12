package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	DBName 		string 	`json:"dbname"`
	DBPort 		string 	`json:"dbport"`
	LogLevel	string	`json:"loglevel"`
	Bindaddr	string	`json:"bindaddr"`
	Cacheaddr	string	`json:"cacheaddr"`
	Cachepass	string	`json:"cachepass"`
}

func ParseConfig(path string) (*Config, error) {
	filebody, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(filebody, &cfg)
	if err != nil {
		return nil, err
	}
	fmt.Println(cfg)
	return &cfg, nil
}
