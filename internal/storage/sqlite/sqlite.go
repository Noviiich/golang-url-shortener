package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Noviiich/golang-url-shortener/internal/storage"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePaths string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePaths)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// SetURL сохраняет URL в базе данных
func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SetURL"

	stmt, err := s.db.Prepare("INSERT INTO urls (url, alias) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		// проверка на уникальность
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	// получение последнего добавленного id
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

// GetURL возвращает оригинальный URL по короткому URL
func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM urls WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var resUrl string

	err = stmt.QueryRow(alias).Scan(&resUrl)
	if err != nil {
		// обработка пустого результата
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resUrl, nil
}

// DeleteURL удаляет URL из базы данных
func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.sqlite.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM urls WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	_, err = stmt.Exec(alias)
	if err != nil {
		// обработка пустого результата
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return nil
}
