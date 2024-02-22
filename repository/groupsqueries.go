package repo

import (
	"log/slog"

	"github.com/michelemendel/dmtmms/entity"
)

func (r *Repo) SelectGroups() ([]entity.Group, error) {
	var groups []entity.Group
	var uuid string
	var name string

	rows, err := r.DB.Query("SELECT uuid, name FROM groups")
	if err != nil {
		slog.Error(err.Error())
		return groups, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&uuid, &name)
		if err != nil {
			slog.Error(err.Error())
			return groups, err
		}
		groups = append(groups, entity.NewGroup(uuid, name))
	}
	err = rows.Err()
	if err != nil {
		slog.Error(err.Error())
		return groups, err
	}
	return groups, nil
}

func (r *Repo) SelectGroup(groupUUID string) (entity.Group, error) {
	var uuid string
	var name string

	err := r.DB.QueryRow("SELECT uuid, name FROM groups WHERE uuid = ?", groupUUID).Scan(&uuid, &name)
	if err != nil {
		slog.Error(err.Error(), "uuid", groupUUID)
		return entity.Group{}, err
	}
	return entity.NewGroup(uuid, name), nil
}

func (r *Repo) DoesGroupNameExist(groupName string) bool {
	var doExist bool
	err := r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM groups WHERE name = ?)", groupName).Scan(&doExist)
	if err != nil {
		slog.Error(err.Error(), "groupname", groupName)
		return false
	}
	return doExist
}

// SELECT g.uuid, g.name FROM groups as g JOIN members_groups as mg on g.uuid = mg.group_uuid WHERE mg.member_uuid = '1df90dea-e0e1-4ed2-8af5-2c475fd52c77';
// Select groups by member
func (r *Repo) SelectGroupsByMember(memberUUID string) ([]entity.Group, error) {
	var groups []entity.Group
	var gUUID string
	var gName string
	rows, err := r.DB.Query(`
	SELECT g.uuid, g.name
	FROM groups as g 
	JOIN members_groups as mg on g.uuid = mg.group_uuid 
	WHERE mg.member_uuid = ?;
	`, memberUUID)
	if err != nil {
		slog.Error(err.Error())
		return groups, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&gUUID, &gName)
		if err != nil {
			slog.Error(err.Error())
			return groups, err
		}
		groups = append(groups, entity.NewGroup(gUUID, gName))
	}
	err = rows.Err()
	if err != nil {
		slog.Error(err.Error())
		return groups, err
	}
	return groups, nil
}

func (r *Repo) SelectGroupUUIDsByMember(memberUUID string) ([]string, error) {
	groups, _ := r.SelectGroupsByMember(memberUUID)
	groupUUIDs := []string{}
	for _, group := range groups {
		groupUUIDs = append(groupUUIDs, group.UUID)
	}
	return groupUUIDs, nil
}
