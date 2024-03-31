package repository

import (
	"App/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type UserPostgres struct {
	conn *sqlx.DB
}

func NewUserPostgres(conn *sqlx.DB) *UserPostgres {
	return &UserPostgres{conn: conn}
}

func (u UserPostgres) Create(in models.UserIn) error {
	stmt := "INSERT INTO users (email, password) VALUES ($1, $2)"

	_, err := u.conn.Exec(stmt, in.Email, in.Password)

	if err != nil {
		log.Err(err).Send()
		return err
	}

	return nil
}

func (u UserPostgres) GetByEmail(email string) (models.User, error) {
	var user models.User

	query := "SELECT id, email, password FROM users WHERE email = $1"

	err := u.conn.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password)

	if err != nil && err.Error() != "sql: no rows in result set" {
		log.Err(err).Send()
		return models.User{}, err
	}

	return user, nil
}

func (u UserPostgres) GetByID(userID int) (models.User, error) {
	var user models.User

	query := "SELECT id, email, password FROM users WHERE id = $1"

	err := u.conn.QueryRow(query, userID).Scan(&user.ID, &user.Email, &user.Password)

	if err != nil {
		log.Err(err).Send()
		return models.User{}, err
	}

	return user, nil
}

func (u *UserPostgres) ChangePassword(userID int, newPassword models.NewPassword) error {
	stmt := "UPDATE users SET password = $1 WHERE id = $2;"

	result, err := u.conn.Exec(stmt, newPassword.Password, userID)

	if err != nil {
		log.Err(err).Send()
		return err
	}

	count, err := result.RowsAffected()

	if err != nil {
		log.Err(err).Send()
		return err
	}

	if count == 0 {
		return NotUpdateRow
	}

	return nil
}
