package repo

import (
	"fmt"
	"log/slog"

	"github.com/michelemendel/dmtmms/e"
	"github.com/michelemendel/dmtmms/entity"
)

func (r *Repo) CreateUser(user entity.User) error {
	result, err := r.DB.Exec("INSERT INTO users(name, password, role) VALUES(?, ?, ?)", user.Name, user.HashedPassword, user.Role)
	fmt.Println("[REPO]:AddUser", "result:", result, "err:", err)
	if err != nil {
		slog.Error(err.Error(), "name", user.Name)
		return e.UserExists
	}
	return nil
}

func (r *Repo) DeleteUser(username string) error {
	result, err := r.DB.Exec("DELETE FROM users WHERE name=?", username)
	fmt.Println("[REPO]:DeleteUser", "result:", result, "err:", err)
	if err != nil {
		slog.Error(err.Error(), "name", username)
		return err
	}
	return nil
}

func (r *Repo) UpdateUser(user entity.User) error {
	_, err := r.DB.Exec("UPDATE users SET role=? WHERE name=?", user.Role, user.Name)
	fmt.Println("[REPO]:UpdateUser", "user:", user, "err:", err)
	if err != nil {
		slog.Error(err.Error(), "name", user.Name)
		return err
	}
	return nil
}

func (r *Repo) UpdateUserPassword(user entity.User) error {
	_, err := r.DB.Exec("UPDATE users SET password=? WHERE name=?", user.HashedPassword, user.Name)
	fmt.Println("[REPO]:UpdateUserPassword", "user:", user, "err:", err)
	if err != nil {
		slog.Error(err.Error(), "name", user.Name)
		return err
	}
	return nil
}

func (r *Repo) UpdateUserRole(username string, role string) error {
	result, err := r.DB.Exec("UPDATE users SET role=? WHERE name=?", role, username)
	fmt.Println("[REPO]:SetUserRole", "result:", result, "err:", err)
	if err != nil {
		slog.Error(err.Error(), "name", username)
		return err
	}
	return nil
}
