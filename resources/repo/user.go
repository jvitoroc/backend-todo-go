package repo

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID           int       `json:"userId" db:"userId"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"passwordHash"`
	Email        string    `json:"email" db:"email"`
	Active       bool      `json:"active" db:"active"`
	CreatedAt    time.Time `json:"createdAt" db:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updatedAt"`
}

type UserActivationRequest struct {
	UserID    int       `json:"userId" db:"userId"`
	Code      string    `json:"code" db:"code"`
	ExpiresAt time.Time `json:"expiresAt" db:"expiresAt"`
}

func InsertUser(db *sqlx.DB, username string, email string, passwordHash string) (*int64, error) {
	now := time.Now()

	res, err := db.NamedExec("insert into user (username, email, passwordHash, createdAt, updatedAt) values (:username, :email, :passwordHash, :createdAt, :updatedAt)",
		map[string]interface{}{
			"username":     username,
			"email":        email,
			"passwordHash": passwordHash,
			"createdAt":    now,
			"updatedAt":    now})

	if err != nil {
		return nil, err
	}

	id, _ := res.LastInsertId()

	return &id, err
}

func UpsertUserActivationRequest(db *sqlx.DB, userId int64, code string, expiresAt time.Time) error {
	_, err := db.NamedExec(`
		insert into user_activation_request (userId, expiresAt, code) values (:userId, :expiresAt, :code)
		on conflict(userId) do update
		set expiresAt = :expiresAt, code = :code
	`,
		map[string]interface{}{
			"userId":    userId,
			"expiresAt": expiresAt,
			"code":      code})

	return err
}

func UpdateUser(db *sqlx.DB, userId int64, pairs map[string]interface{}) (sql.Result, error) {
	var command string = "update user set "

	for key := range pairs {
		command = command + key + " = :" + key + ","
	}

	command = command[:len(command)-1]
	command = command + " where userId = :userId"

	pairs["userId"] = userId
	pairs["updatedAt"] = time.Now()

	res, err := db.NamedExec(command, pairs)

	return res, err
}

func GetUserActivationRequest(db *sqlx.DB, userId int64) (*UserActivationRequest, error) {
	request := UserActivationRequest{}

	err := db.Get(&request, `select * from user_activation_request where userId = $1`, userId)

	return &request, err
}

func GetUser(db *sqlx.DB, userId int64) (*User, error) {
	user := User{}

	err := db.Get(&user, `select * from user where userId = $1`, userId)

	return &user, err
}

func GetUserByUsername(db *sqlx.DB, username string) (*User, error) {
	user := User{}

	err := db.Get(&user, `select * from user where username = $1`, username)

	return &user, err
}
