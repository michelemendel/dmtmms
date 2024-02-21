package repo

import (
	"log/slog"

	"github.com/michelemendel/dmtmms/entity"
)

func (r *Repo) SelectFamilies() ([]entity.Family, error) {
	var families []entity.Family
	var uuid string
	var name string

	rows, err := r.DB.Query("SELECT uuid, name FROM families")
	if err != nil {
		slog.Error(err.Error())
		return families, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&uuid, &name)
		if err != nil {
			slog.Error(err.Error())
			return families, err
		}
		if uuid != "0" {
			families = append(families, entity.NewFamily(uuid, name))
		}
	}
	err = rows.Err()
	if err != nil {
		slog.Error(err.Error())
		return families, err
	}
	return families, nil
}

func (r *Repo) SelectFamily(familyUUID string) (entity.Family, error) {
	var uuid string
	var name string

	err := r.DB.QueryRow("SELECT uuid, name FROM families WHERE uuid = ?", familyUUID).Scan(&uuid, &name)
	if err != nil {
		slog.Error(err.Error(), "uuid", familyUUID)
		return entity.Family{}, err
	}
	return entity.NewFamily(uuid, name), nil
}

func (r *Repo) GetFamilyNameByUUID(familyUUID string) (string, error) {
	var name string
	err := r.DB.QueryRow("SELECT name FROM families WHERE uuid=?", familyUUID).Scan(&name)
	if err != nil {
		slog.Error(err.Error(), "uuid", familyUUID)
		return "", err
	}
	return name, nil
}
