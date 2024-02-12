package repo

import (
	"log/slog"

	"github.com/michelemendel/dmtmms/entity"
)

func (r *Repo) SelectUsers() ([]entity.User, error) {
	var users []entity.User
	var name string
	var role string

	rows, err := r.DB.Query("SELECT name, role FROM users")
	if err != nil {
		slog.Error(err.Error())
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name, &role)
		if err != nil {
			slog.Error(err.Error())
			return users, err
		}
		users = append(users, entity.NewUser(name, "", role))
	}
	err = rows.Err()
	if err != nil {
		slog.Error(err.Error())
		return users, err
	}
	return users, nil
}

func (r *Repo) SelectUser(username string) (entity.User, error) {
	var name string
	var pw string
	var role string

	err := r.DB.QueryRow("SELECT name, password, role FROM users WHERE name = ?", username).Scan(&name, &pw, &role)
	if err != nil {
		slog.Error(err.Error(), "name", username)
		return entity.User{}, err
	}
	return entity.NewUser(name, pw, role), nil
}

func (r *Repo) DoesUsernameExist(username string) bool {
	var name string
	err := r.DB.QueryRow("SELECT name FROM users WHERE name = ?", username).Scan(&name)
	if err != nil {
		slog.Error(err.Error(), "name", username)
		return false
	}
	return true
}
