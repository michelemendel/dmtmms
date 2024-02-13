package repo

import (
	"log/slog"

	"github.com/michelemendel/dmtmms/e"
	"github.com/michelemendel/dmtmms/entity"
)

func (r *Repo) CreateGroup(group entity.Group) error {
	_, err := r.DB.Exec("INSERT INTO groups(uuid, name) VALUES(?, ?)", group.UUID, group.Name)
	if err != nil {
		slog.Error(err.Error(), "uuid", group.UUID, "name", group.Name)
		return e.ErrGroupExists
	}
	slog.Info("CreateGroup", "uuid", group.UUID, "name", group.Name)
	return nil
}

func (r *Repo) DeleteGroup(groupUUID string) error {
	_, err := r.DB.Exec("DELETE FROM groups WHERE uuid=?", groupUUID)
	if err != nil {
		slog.Error(err.Error(), "uuid", groupUUID)
		return e.ErrGroupIsUsedByMembers
	}
	slog.Info("DeleteGroup", "uuid", groupUUID)
	return nil
}

func (r *Repo) UpdateGroup(group entity.Group) error {
	_, err := r.DB.Exec("UPDATE groups SET name=? WHERE uuid=?", group.Name, group.UUID)
	if err != nil {
		slog.Error(err.Error(), "uuid", group.UUID, "name", group.Name)
		return err
	}
	slog.Info("UpdateGroup", "uuid", group.UUID, "name", group.Name)
	return nil
}
