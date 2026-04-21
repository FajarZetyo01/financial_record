package models

import (
	"database/sql"
	"financial_record/entities"
	"time"
)

type UserModel struct {
	db *sql.DB
}

func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{
		db: db,
	}
}

func (model UserModel) FindUserById(id string) (entities.User, error) {

	var user entities.User
	query := "SELECT email, name, photo FROM users WHERE id = ?"

	err := model.db.QueryRow(query, id).Scan(&user.Email, &user.Name, &user.Photo)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (model UserModel) UpdateProfile(user entities.User) error {

	var query string
	var args []interface{}

	if user.Password != "" {
		query = "UPDATE users SET name = ?, password = ?, photo = ?, updated_at = ? WHERE id = ?"
		args = []interface{}{user.Name, user.Password, user.Photo, time.Now(), user.Id}
	} else {
		query = "UPDATE users SET name = ?, photo = ?, updated_at = ? WHERE id = ?"
		args = []interface{}{user.Name, user.Photo, time.Now(), user.Id}
	}

	_, err := model.db.Exec(query, args...)
	return err
}
