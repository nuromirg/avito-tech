package store

import (
	"avito_task/config"
	"database/sql"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type Storage struct {
	Db  *sql.DB
	r   *redis.Client
	cfg *config.Config
	u   *UserRepository
}

func CreateStorage(_cfg config.Config) *Storage {
	return &Storage{
		cfg: &_cfg,
	}
}

func (s *Storage) Open() error {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:" + s.cfg.DBPort + ")/"+s.cfg.DBName)
	if err != nil {
		logrus.Printf("Error %s when creating DB\n", err)
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	s.Db = db

	rClient := redis.NewClient(&redis.Options{
		Addr: 		s.cfg.Cacheaddr,
		Password: 	s.cfg.Cachepass,
		DB: 		0,
	})
	s.r = rClient
	ping, err := s.r.Ping(s.r.Context()).Result()
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Println(ping)

	return nil
}

func (s *Storage) Close() {
	err := s.Db.Close()
	if err != nil {
		return 
	}
}

func (s *Storage) User() *UserRepository {
	if s.u != nil {
		return s.u
	}
	s.u = &UserRepository{
		Storage: s,
	}
	return s.u
}
