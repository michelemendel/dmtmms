package repo

import (
	"log/slog"

	// "github.com/michelemendel/dmtmms/e"
	"github.com/michelemendel/dmtmms/entity"
)

func (r *Repo) CreateMember(member entity.Member) error {
	// _, err := r.DB.Exec("INSERT INTO groups(uuid, name) VALUES(?, ?)", group.UUID, group.Name)
	// if err != nil {
	// slog.Error(err.Error(), "uuid", group.UUID, "name", group.Name)
	// return e.ErrGroupExists
	// }
	// slog.Info("CreateGroup", "uuid", group.UUID, "name", group.Name)
	return nil
}

func (r *Repo) DeleteMember(memberUUID string) error {
	_, err := r.DB.Exec("DELETE FROM members WHERE uuid=?", memberUUID)
	if err != nil {
		slog.Error(err.Error(), "uuid", memberUUID)
		return err
	}
	slog.Info("DeleteMember", "uuid", memberUUID)
	return nil
}

func (r *Repo) UpdateMember(member entity.Member) error {
	// _, err := r.DB.Exec("UPDATE members SET name=? WHERE uuid=?", member.Name, member.UUID)
	// if err != nil {
	// slog.Error(err.Error(), "uuid", member.UUID, "name", member.Name)
	// return err
	// }
	slog.Info("UpdateGroup", "uuid", member.UUID, "name", member.Name)
	return nil
}
