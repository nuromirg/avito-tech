package store_test

import (
	"avito_task/internal/model"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"os"
	"testing"
)

var (
	databaseURL string
)

var u = &model.User{
	Id: 13,
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "host=localhost dbname=avito_test sslmode=disable"
	}

	os.Exit(m.Run())
}

/*
func TestUserRepository_FindById(t *testing.T) {
	db, mock := NewMock()
	st := store.Storage{
		Db: db,

	}
	repo := &store.UserRepository{ }
}
*/