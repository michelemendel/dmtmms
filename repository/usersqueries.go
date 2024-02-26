package repo

import (
	"log/slog"

	"github.com/michelemendel/dmtmms/entity"
)

func (r *Repo) SelectUsers() ([]entity.User, error) {
	var users []entity.User

	rows, err := r.DB.Query("SELECT name, role FROM users")
	if err != nil {
		slog.Error(err.Error())
		return users, err
	}
	defer rows.Close()
	var u entity.User
	for rows.Next() {
		err := rows.Scan(&u.Name, &u.Role)
		if err != nil {
			slog.Error(err.Error())
			return users, err
		}
		users = append(users, entity.NewUser(u.Name, "", u.Role))
	}
	err = rows.Err()
	if err != nil {
		slog.Error(err.Error())
		return users, err
	}
	return users, nil
}

func (r *Repo) SelectUser(username string) (entity.User, error) {
	var u entity.User
	err := r.DB.QueryRow("SELECT name, password, role FROM users WHERE name = ?", username).Scan(&u.Name, &u.HashedPassword, &u.Role)
	if err != nil {
		slog.Error(err.Error(), "name", username)
		return entity.User{}, err
	}
	return entity.NewUser(u.Name, u.HashedPassword, u.Role), nil
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
