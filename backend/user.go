package backend

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"os"
	"time"

	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
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

func generateRandomString() (string, error) {
	b := make([]byte, 30)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func generateToken(userID int, validFor int) (*User, error) {
	var id int
	var user User
	token, err := generateRandomString()
	expiry := time.Now().UTC().Unix() + int64(validFor)
	if err != nil {
		return nil, err
	}
	err = db.QueryRow(`INSERT INTO token(token, user_id, expiry_date)
		VALUES ($1, $2, $3)
		RETURNING token_id`, token, userID, expiry).Scan(&id)
	if err != nil {
		return nil, err
	}
	err = db.QueryRow(`SELECT users.user_id, users.email, users.role, token.token
		FROM token
		INNER JOIN users
		ON token.user_id = users.user_id
		WHERE token_id = $1`, id).Scan(&user.ID, &user.Email, &user.Role, &user.Token)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func sendTokenMail(email string, token string) {
	from := mail.NewEmail("Peak Power Leasing", "mail@peakpowerleasing.de")
	subject := "Your Peak Power Leasing Token"
	to := mail.NewEmail(email, email)
	content := mail.NewContent("text/plain", "Your temporary Peak Power Leasing Login token is:\n\n"+token)
	m := mail.NewV3MailInit(from, subject, to, content)

	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	_, err := sendgrid.API(request)
	if err != nil {
		log.Println("Error while sending mail:", err)
	}
}

func LoginPwd(email string, passwd string) (*User, error) {
	var id int
	var hash string
	err := db.QueryRow(`SELECT users.user_id, password.password
		FROM users
		INNER JOIN password
		ON users.user_id = password.user_id
		WHERE users.email = $1 AND users.role IN ('staff', 'owner')`, email).Scan(&id, &hash)
	if err != nil {
		return nil, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd)); err != nil {
		return nil, errors.New("Wrong password")
	}
	return generateToken(id, 7*24*60*60)
}

func RequestToken(email string) (bool, error) {
	var id int
	err := db.QueryRow(`SELECT users.user_id
		FROM users
		WHERE users.email = $1`, email).Scan(&id)
	if err != nil {
		return false, err
	}
	user, err := generateToken(id, 30*60)
	if err == nil {
		sendTokenMail(user.Email, user.Token)
		return true, err
	} else {
		return false, err
	}
}

func ValidateToken(email string, token string) (*User, error) {
	var user User
	var expiryDate int64
	err := db.QueryRow(`SELECT users.user_id, users.email, users.role, token.token, token.expiry_date
		FROM token
		INNER JOIN users
		ON token.user_id = users.user_id
		WHERE  token.token = $1`, token).Scan(&user.ID, &user.Email, &user.Role, &user.Token, &expiryDate)
	if err != nil {
		return nil, err
	}
	if time.Now().UTC().Unix() >= expiryDate {
		return nil, errors.New("Token expired")
	}
	return &user, nil
}
