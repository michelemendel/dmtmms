package repo

import (
	"log/slog"

	"github.com/michelemendel/dmtmms/entity"
)

func (r *Repo) SelectFamilies() ([]entity.Family, error) {
	var families []entity.Family
	var f entity.Family
	rows, err := r.DB.Query("SELECT uuid, name FROM families")
	if err != nil {
		slog.Error(err.Error())
		return families, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&f.UUID, &f.Name)
		if err != nil {
			slog.Error(err.Error())
			return families, err
		}
		families = append(families, entity.NewFamily(f.UUID, f.Name))
	}
	err = rows.Err()
	if err != nil {
		slog.Error(err.Error())
		return families, err
	}
	// Remove family "none"
	for i, f := range families {
		if f.Name == "none" {
			families = append(families[:i], families[i+1:]...)
			break
		}
	}

	return families, nil
}

func (r *Repo) SelectFamily(familyUUID string) (entity.Family, error) {
	var f entity.Family
	err := r.DB.QueryRow("SELECT uuid, name FROM families WHERE uuid = ?", familyUUID).Scan(&f.UUID, &f.Name)
	if err != nil {
		slog.Error(err.Error(), "uuid", familyUUID)
		return entity.Family{}, err
	}
	return entity.NewFamily(f.UUID, f.Name), nil
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
