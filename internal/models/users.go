package models

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
}

func (m *UserModel) Insert(name, email, password string) error {
	// Step 1: Hash the plain password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	// Step 2: Prepare the INSERT SQL statement
	stmt := `INSERT INTO users (name, email, hashed_password, created)
             VALUES (?, ?, ?, UTC_TIMESTAMP())`

	// Step 3: Try to insert into the database
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		// Step 4: If there's an error, check if it’s a MySQL error
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			// Step 5: If error is 1062 and relates to the UNIQUE email constraint, return a custom error
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		// Step 6: If it's some other error, return it as-is
		return err
	}

	// Success!
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
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

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool
	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"
	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}
