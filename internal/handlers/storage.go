package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
)

type SubmissionStore interface {
	Save(context.Context, Submission) error
}

type FileSubmissionStore struct {
	path string
}

type PostgresSubmissionStore struct {
	db *sql.DB
}

func NewSubmissionStore(submissionsPath string) (SubmissionStore, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		db, err := sql.Open("postgres", databaseURL)
		if err != nil {
			return nil, fmt.Errorf("open postgres connection: %w", err)
		}

		db.SetConnMaxLifetime(5 * time.Minute)
		db.SetMaxIdleConns(2)
		db.SetMaxOpenConns(10)

		if err := db.Ping(); err != nil {
			return nil, fmt.Errorf("ping postgres: %w", err)
		}

		return &PostgresSubmissionStore{db: db}, nil
	}

	return &FileSubmissionStore{path: submissionsPath}, nil
}

func (s *FileSubmissionStore) Save(_ context.Context, submission Submission) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}

	file, err := os.OpenFile(s.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(submission)
}

func (s *PostgresSubmissionStore) Save(ctx context.Context, submission Submission) error {
	const query = `
		INSERT INTO contact_submissions (name, email, message, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := s.db.ExecContext(ctx, query, submission.Name, submission.Email, submission.Message, submission.CreatedAt)
	return err
}
