//package models
//
//import (
//	"database/sql"
//	"financial_record/entities"
//)
//
//type AuthModel struct {
//	db *sql.DB // AuthController tidak mengacu ke entities/auth, tapi ke database connection (*sql.DB).
//}
//
//func NewAuthModel(db *sql.DB) *AuthModel {
//	return &AuthModel{
//		db: db,
//	}
//}
//func (model AuthModel) Register(register entities.Register) error {
//	_, err := model.db.Exec("INSERT INTO users (id,name,email,password) VALUES (?,?,?,?)",
//		register.Id, register.Name, register.Email, register.Password,
//	)
//	return err
//}
//func (model AuthModel) Login(email string) (entities.Auth, error) {
//	var auth entities.Auth
//	query := "SELECT id, name, email, password FROM users WHERE email = ?"
//	err := model.db.QueryRow(query, email).Scan(
//		&auth.Id,
//		&auth.Name,
//		&auth.Email,
//		&auth.Password,
//	)
//	if err != nil {
//		return auth, nil
//	}
//	return auth, nil
//}

package models

import (
	"database/sql"
	"errors"
	"financial_record/entities"
)

type AuthModel struct {
	db *sql.DB
}

func NewAuthModel(db *sql.DB) *AuthModel {
	return &AuthModel{
		db: db,
	}
}

func (m *AuthModel) Register(register entities.Register) error {
	_, err := m.db.Exec("INSERT INTO users (id, name, email, password) VALUES (?, ?, ?, ?)",
		register.Id,
		register.Name,
		register.Email,
		register.Password,
	)
	return err
}

func (m *AuthModel) Login(email string) (entities.Auth, error) {
	var auth entities.Auth
	query := `SELECT id, name, email, password FROM users WHERE email = ?`
	err := m.db.QueryRow(query, email).Scan(
		&auth.Id,
		&auth.Name,
		&auth.Email,
		&auth.Password,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auth, err // user tidak ditemukan
		}
		return auth, err // error DB
	}
	return auth, nil
}
