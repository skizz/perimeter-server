package main

import (
	"database/sql"
	"time"
)

type Session struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Session) getSession(db *sql.DB) error {
	return db.QueryRow("SELECT id, created_at FROM sessions WHERE id=$1",
		&s.Id).Scan(&s.CreatedAt)
}

func getSessions(db *sql.DB) ([]Session, error) {
	rows, err := db.Query(
		"SELECT id, created_at FROM sessions")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sessions := []Session{}

	for rows.Next() {
		var s Session
		if err := rows.Scan(&s.Id, &s.CreatedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}

	return sessions, nil
}

func (s *Session) createSession(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO sessions DEFAULT VALUES RETURNING id, created_at").Scan(&s.Id, &s.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
