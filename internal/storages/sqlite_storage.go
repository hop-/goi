package storages

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteStorage struct {
	db *sql.DB
}

func (s *SqliteStorage) Close() {
	s.db.Close()
}

// TODO: add Storage implementation

func newSqliteStorage(filename string) (Storage, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}

	return &SqliteStorage{
		db: db,
	}, nil
}
