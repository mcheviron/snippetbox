package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type Snippets struct {
	DB *sql.DB
}

func (m *Snippets) Insert(title, content string, expires int) (int, error) {
	result, err := m.DB.Exec(
		`
                         INSERT INTO
                           snippets (title, content, created, expires)
                         VALUES
                           (
                             ?,
                             ?,
                             UTC_TIMESTAMP (),
                             DATE_ADD (UTC_TIMESTAMP (), INTERVAL ? DAY)
                           )
	`, title, content, expires,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *Snippets) Get(id int) (*Snippet, error) {
	s := &Snippet{}
	err := m.DB.QueryRow(
		`
                     SELECT
                       id,
                       title,
                       content,
                       expires
                     FROM
                       snippets
                     WHERE
                       expires > UTC_TIMESTAMP ()
                       AND id = ?
	`, id,
	).Scan(&s.ID, &s.Title, &s.Content, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

// * This will return the last 10 snippets added
func (m *Snippets) Latest() ([]*Snippet, error) {
	rows, err := m.DB.Query(
		`
                        SELECT
                          id,
                          title,
                          content,
                          created,
                          expires
                        FROM
                          snippets
                        WHERE
                          expires > UTC_TIMESTAMP ()
                        ORDER BY
                          id DESC
                        LIMIT
                          10
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	snippets := []*Snippet{}

	for rows.Next() {
		s := Snippet{}
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
