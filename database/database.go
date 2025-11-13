package database

import (
	"context"
	"database/sql"

	"github.com/amjadnzr/url-shortly/models"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	*sql.DB
}

func InitDatabase() (*Database, error) {
	db, err := sql.Open("sqllite3", "url-shortly.db")
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &Database{db}, nil
}

func (db *Database) CreateNewUser(c context.Context, user *models.User) (int64, error) {
	row, err := db.ExecContext(c,
		"INSERT INTO Users (name, email, password) VALUES (?, ?, ?)",
		user.Name, user.Email, user.PasswordHash,
	)
	if err != nil {
		return 0, err
	}
	id, err := row.LastInsertId()
	if err != nil {
		return id, err
	}

	return id, err
}

func (db *Database) GetUserById(c context.Context, id int64) (*models.User, error) {
	user := new(models.User)

	if err :=
		db.QueryRow("SELECT * FROM Users where id = ?", id).
			Scan(user); err != nil {
		return nil, err
	}
	return user, nil
}
