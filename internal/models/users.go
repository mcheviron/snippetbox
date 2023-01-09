package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// UserModel
type Users struct {
	DB *sql.DB
}

func (u *Users) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	_, err = u.DB.Exec(`
                   INSERT INTO
                     users (name, email, hashed_password, created)
                   VALUES
                     (?, ?, ?, UTC_TIMESTAMP ())
	`, name, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 &&
				strings.Contains(mySQLError.Message, "user_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (u *Users) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	err := u.DB.QueryRow(`
                     SELECT
                       id,
                       hashed_password
                     FROM
                       users
                     WHERE
                       email = ?
	`, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

func (u *Users) Exists(id int) (bool, error) {
	var exists bool
	err := u.DB.QueryRow(`
                     SELECT
                       EXISTS (
                         SELECT
                           TRUE
                         FROM
                           users
                         WHERE
                           id = ?
                       )
	`, id).Scan(&exists)
	return exists, err
}
