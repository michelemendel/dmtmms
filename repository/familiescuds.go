package repo

import (
	"log/slog"

	"github.com/michelemendel/dmtmms/e"
	"github.com/michelemendel/dmtmms/entity"
)

func (r *Repo) CreateFamily(family entity.Family) error {
	_, err := r.DB.Exec("INSERT INTO families(uuid, name) VALUES(?, ?)", family.UUID, family.Name)
	if err != nil {
		slog.Error(err.Error(), "uuid", family.UUID, "name", family.Name)
		return e.ErrFamilyExists
	}
	slog.Info("CreateFamily", "uuid", family.UUID, "name", family.Name)
	return nil
}

func (r *Repo) DeleteFamily(familyUUID string) error {
	_, err := r.DB.Exec("DELETE FROM families WHERE uuid=?", familyUUID)
	if err != nil {
		slog.Error(err.Error(), "uuid", familyUUID)
		return e.ErrFamilyIsUsedByMembers
	}
	slog.Info("DeleteFamily", "uuid", familyUUID)
	return nil
}

func (r *Repo) UpdateFamily(family entity.Family) error {
	_, err := r.DB.Exec("UPDATE families SET name=? WHERE uuid=?", family.Name, family.UUID)
	if err != nil {
		slog.Error(err.Error(), "uuid", family.UUID, "name", family.Name)
		return err
	}
	slog.Info("UpdateFamily", "uuid", family.UUID, "name", family.Name)
	return nil
}
