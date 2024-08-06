package storages

import (
	"database/sql"
	"fmt"

	"github.com/hop-/goi/internal/core"
	"github.com/hop-/golog"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteStorage struct {
	db *sql.DB
}

func checkForeignKeys(s *SqliteStorage) (bool, error) {
	fkCheck := s.db.QueryRow("PRAGMA foreign_keys")
	if fkCheck == nil {
		return false, fmt.Errorf("failed to get foreign keys pragma")
	}
	err := fkCheck.Err()
	if err != nil {
		return false, err
	}

	var v int
	fkCheck.Scan(&v)

	return v == 1, nil
}

func sqliteChecksAndSetup(s *SqliteStorage) error {
	// Check foreign_keys
	fkStatus, err := checkForeignKeys(s)
	if err != nil {
		return err
	}

	if !fkStatus {
		// return fmt.Errorf("foreign keys feature is disabled")
		_, err = s.db.Exec("PRAGMA foreign_keys = ON")
		if err != nil {
			return err
		}
		// Check foreign keys once more
		fkStatus, err = checkForeignKeys(s)
		if err != nil {
			return err
		} else if !fkStatus {
			return fmt.Errorf("failed to enable foreign keys")
		}
	}

	return nil
}

func createTopicsTable(s *SqliteStorage) error {
	// Create table
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS topics (
		id INTEGER PRIMARY KEY,
		name TEXT
	)`)
	if err != nil {
		return err
	}

	// Create index
	_, err = s.db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS uq_topics_name ON topics (name)`)
	return err
}

func createConsumerGroupsTable(s *SqliteStorage) error {
	// Create table
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS consumer_groups (
		id INTEGER PRIMARY KEY,
		name TEXT
	)`)
	if err != nil {
		return err
	}

	// Create index
	_, err = s.db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS uq_consumer_groups_name ON consumer_groups (name)`)
	return err
}

func createConsumerGroupTopicStates(s *SqliteStorage) error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS consumer_group_topic_states (
		topic_id INTEGER NOT NULL,
		consumer_group_id INTEGER NOT NULL,
		PRIMARY KEY (topic_id, consumer_group_id),
		FOREIGN KEY (topic_id) REFERENCES topics(id),
		FOREIGN KEY (consumer_group_id) REFERENCES consumer_groups(id)
	)`)
	return err
}

func createMessagesTable(s *SqliteStorage) error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY,
		occurred_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		content BLOB,
		topic_id INTEGER NOT NULL,
		FOREIGN KEY (topic_id) REFERENCES topics(id)
	)`)
	if err != nil {
		return err
	}

	// Create indeces
	_, err = s.db.Exec(`CREATE INDEX IF NOT EXISTS ix_message_occured_at ON messages (occurred_at DESC)`)
	return err
}

func (s *SqliteStorage) Init() error {
	golog.Debug("Setup sqlite storage")
	err := sqliteChecksAndSetup(s)
	if err != nil {
		return err
	}

	golog.Debug("Creating consumers groups table")
	err = createConsumerGroupsTable(s)
	if err != nil {
		return err
	}

	golog.Debug("Creating topics table")
	err = createTopicsTable(s)
	if err != nil {
		return err
	}

	golog.Debug("Creating consumer group topic states table")
	err = createConsumerGroupTopicStates(s)
	if err != nil {
		return err
	}

	golog.Debug("Creating messages table")
	return createMessagesTable(s)
}

func (s *SqliteStorage) Close() error {
	return s.db.Close()
}

func (s *SqliteStorage) Topics() ([]core.Topic, error) {
	// TODO
	return []core.Topic{}, nil
}

func (s *SqliteStorage) NewTopic(*core.Topic) error {
	// TODO
	return nil
}

func (s *SqliteStorage) ConsumerGroups() ([]core.ConsumerGroup, error) {
	// TODO
	return []core.ConsumerGroup{}, nil
}

func (s *SqliteStorage) NewConsumerGroup(*core.ConsumerGroup) error {
	// TODO
	return nil
}

func (s *SqliteStorage) Messages(*core.Topic) ([]core.Message, error) {
	// TODO
	return []core.Message{}, nil
}

func (s *SqliteStorage) NewMessage(*core.Message) error {
	// TODO
	return nil
}

// TODO: add Storage implementation

func newSqliteStorage(filename string) (core.Storage, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}

	return &SqliteStorage{
		db: db,
	}, nil
}
