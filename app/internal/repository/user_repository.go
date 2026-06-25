package repository

import (
	"github.com/Ozielmarcos/mytodolist/app/internal/model"
	"github.com/Ozielmarcos/mytodolist/app/pkg/database"
)

func CreateUser(user model.User) error {
	_, err := database.DB.Exec(
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3)",
		user.Name, user.Email, user.Password,
	)

	if err != nil {
		return err
	}
	return nil
}

func GetUserByEmail(email string) (model.User, error) {
	var user model.User

	err := database.DB.QueryRow(
		"SELECT * FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func UpdateUser(id string, user model.User) error {
	sql := "UPDATE users SET name = $1, email = $2, password = $3 WHERE id=$5"

	_, err := database.DB.Exec(sql, user.Name, user.Email, user.Password, id)

	if err != nil {
		return err
	}

	return nil
}
