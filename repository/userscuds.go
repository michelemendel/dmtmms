package repo

import (
	"log/slog"

	"github.com/michelemendel/dmtmms/e"
	"github.com/michelemendel/dmtmms/entity"
)

func (r *Repo) CreateUser(user entity.User) error {
	_, err := r.DB.Exec("INSERT INTO users(name, password, role) VALUES(?, ?, ?)", user.Name, user.HashedPassword, user.Role)
	if err != nil {
		slog.Error(err.Error(), "name", user.Name)
		return e.ErrUserExists
	}
	slog.Info("CreateUser", "name", user.Name, "role", user.Role)
	return nil
}

func (r *Repo) DeleteUser(username string) error {
	_, err := r.DB.Exec("DELETE FROM users WHERE name=?", username)
	if err != nil {
		slog.Error(err.Error(), "name", username)
		return err
	}
	slog.Info("DeleteUser", "name", username)
	return nil
}

func (r *Repo) UpdateUser(user entity.User) error {
	_, err := r.DB.Exec("UPDATE users SET role=? WHERE name=?", user.Role, user.Name)
	if err != nil {
		slog.Error(err.Error(), "name", user.Name)
		return err
	}
	slog.Info("UpdateUser", "name", user.Name, "role", user.Role)
	return nil
}

func (r *Repo) UpdateUserPassword(user entity.User) error {
	_, err := r.DB.Exec("UPDATE users SET password=? WHERE name=?", user.HashedPassword, user.Name)
	if err != nil {
		slog.Error(err.Error(), "name", user.Name)
		return err
	}
	slog.Info("UpdateUserPassword", "name", user.Name)
	return nil
}

func (r *Repo) UpdateUserRole(username string, role string) error {
	_, err := r.DB.Exec("UPDATE users SET role=? WHERE name=?", role, username)
	if err != nil {
		slog.Error(err.Error(), "name", username)
		return err
	}
	slog.Info("SetUserRole", "name", username, "role", role)
	return nil
}

func (r *Repo) ResetPassword(username string, hashedPassword string) error {
	_, err := r.DB.Exec("UPDATE users SET password = ? WHERE name = ?", hashedPassword, username)
	if err != nil {
		slog.Error(err.Error(), "name", username)
		return err
	}
	slog.Info("ResetPassword", "name", username)
	return nil
}
