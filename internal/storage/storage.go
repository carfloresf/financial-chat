package storage

import (
	"time"

	"database/sql"

	log "github.com/sirupsen/logrus"

	"github.com/carfloresf/financial-chat/config"
	"github.com/carfloresf/financial-chat/internal/model"
)

type DB struct {
	Conn *sql.DB
}

type Storage interface {
	CreateUser(username, passwordHash string) (int, error)
	GetUser(username string) (*model.User, error)
}

func NewStorage(config *config.DB) (*DB, error) {
	sqliteDatabase, err := sql.Open("sqlite3", config.DBFile)
	if err != nil {
		log.Fatal("error opening Conn connection: %w", err)

		return nil, err
	}

	db := DB{
		Conn: sqliteDatabase,
	}

	return &db, nil
}

func (s *DB) Close() error {
	return s.Conn.Close()
}

const (
	insertUserQuery = `INSERT
						INTO users (username, passwordHash, created_at)
						VALUES    ($1, $2, $3)`
	getUserQuery = `SELECT id, username, passwordHash, created_at
					FROM users
					WHERE username = $1`
)

func (s *DB) CreateUser(username, passwordHash string) (int, error) {
	res, err := s.Conn.Exec(insertUserQuery, username, passwordHash, time.Now())
	if err != nil {
		return 0, err
	}

	var generatedID int64

	if generatedID, err = res.LastInsertId(); err != nil {
		return 0, err
	}

	return int(generatedID), nil
}

func (s *DB) GetUser(username string) (*model.User, error) {
	var user model.User

	err := s.Conn.QueryRow(getUserQuery, username).
		Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
