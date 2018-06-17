package backend

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Token string `json:"token"`
}

func InsertUser(user *User) error {
	var id int
	err := db.QueryRow(`
		INSERT INTO users(email)
		VALUES ($1)
		RETURNING id
	`, user.Email).Scan(&id)
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

func GetUserByID(id int) (*User, error) {
	var email string
	err := db.QueryRow("SELECT email FROM users WHERE user_id=$1", id).Scan(&email)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:    id,
		Email: email,
	}, nil
}

func RemoveUserByID(id int) error {
	_, err := db.Exec("DELETE FROM users WHERE id=$1", id)
	return err
}

func generateToken(userID int) (*User, error) {
	return nil, errors.New("generateToken not implemented")
}

func LoginPwd(email string, passwd string) (*User, error) {
	var id int
	var hash string
	err := db.QueryRow(`SELECT users.user_id, password.password
		FROM users
		INNER JOIN password
		ON users.user_id = password.user_id
		WHERE users.email=$1 AND users.role IN ('staff', 'owner')`, email).Scan(&id, &hash)
	if err != nil {
		return nil, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd)); err != nil {
		return nil, errors.New("Wrong password")
	}
	return generateToken(id)
}
