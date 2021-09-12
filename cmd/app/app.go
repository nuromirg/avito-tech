package main

import (
	"avito_task/config"
	"avito_task/internal/server"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.ParseConfig("config/config.json")
	if err != nil {
		logrus.Error(err)
	}

	//db := store.CreateStorage(*cfg)
	s := server.New(cfg)
	if err := s.Start(); err != nil {
		logrus.Error(err)
	}

}
