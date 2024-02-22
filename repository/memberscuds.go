package repo

import (
	// "fmt"

	"log/slog"

	// "github.com/michelemendel/dmtmms/e"
	"github.com/michelemendel/dmtmms/e"
	"github.com/michelemendel/dmtmms/entity"
)

func (r *Repo) CreateMember(member entity.Member, groupUUIDs []string) error {
	tx, _ := r.DB.Begin()
	familyUUID := "0"
	familyName := "none"
	if member.FamilyUUID != "" {
		familyUUID = member.FamilyUUID
		familyName = member.FamilyName
	}

	_, err := tx.Exec(`
	INSERT INTO members(
		uuid, 
		id, name, dob, personnummer, email, mobile, 
		address1, address2, postnummer, poststed, 
		synagogue_seat, membership_fee_tier, registered_date, deregistered_date, 
		receive_email, receive_mail, receive_hatikva, archived, status, 
		family_uuid, family_name
		) VALUES(
			?, 
			?, ?, julianday(?), ?, ?, ?, 
			?, ?, ?, ?, 
			?, ?, julianday(?), julianday(?), 
			?, ?, ?, ?, ?, 
			?, ?
		)
		`,
		member.UUID,
		member.ID, member.Name, member.DOB, member.Personnummer, member.Email, member.Mobile,
		member.Address.Address1, member.Address.Address2, member.Address.Postnummer, member.Address.Poststed,
		member.Synagogueseat, member.MembershipFeeTier, member.RegisteredDate, member.DeregisteredDate,
		member.ReceiveEmail, member.ReceiveMail, member.ReceiveHatikva, member.Archived, member.Status,
		familyUUID, familyName,
	)
	if err != nil {
		slog.Error(err.Error(), "uuid", member.UUID, "name", member.Name, "familyUUID", familyUUID, "familyName", familyName)
		tx.Rollback()
		return e.ErrCreatingMember
	}

	// Add member to groups
	groupUUIDs = pruneGroupUUIDs(groupUUIDs)
	for _, groupUUID := range groupUUIDs {
		_, err = tx.Exec(`INSERT INTO members_groups(member_uuid, group_uuid) VALUES(?, ?)`, member.UUID, groupUUID)
		if err != nil {
			slog.Error(err.Error(), "uuid", member.UUID, "groupUUD", groupUUID)
			tx.Rollback()
			return e.ErrAddingGroupForMember
		}
	}

	tx.Commit()
	slog.Info("CreateMember", "uuid", member.UUID, "name", member.Name)
	return nil
}

func (r *Repo) ArchiveMember(memberUUID string) error {
	_, err := r.DB.Exec("UPDATE members SET archived=true WHERE uuid=?", memberUUID)
	if err != nil {
		slog.Error(err.Error(), "uuid", memberUUID)
		return err
	}
	slog.Info("ArchivedMember", "uuid", memberUUID)
	return nil
}

func (r *Repo) DeleteMember(memberUUID string) error {
	_, err := r.DB.Exec("DELETE FROM members_groups WHERE member_uuid=?", memberUUID)
	if err != nil {
		slog.Error("delete from mebers_groups", "error", err.Error(), "uuid", memberUUID)
		return err
	}

	_, err = r.DB.Exec("DELETE FROM members WHERE uuid=?", memberUUID)
	if err != nil {
		slog.Error("delete from members", "error", err.Error(), "uuid", memberUUID)
		return err
	}
	slog.Info("DeleteMember", "uuid", memberUUID)
	return nil
}

// slog.Info("UpdateGroup", "uuid", member.UUID, "name", member.Name)
// SELECT m.name,m.family_name,mg.member_uuid,mg.group_uuid,g.name from members as m LEFT JOIN members_groups as mg on m.uuid=mg.member_uuid LEFT JOIN groups as g ON group_uuid=g.uuid;
func (r *Repo) UpdateMember(member entity.Member, groupUUIDs []string) error {
	tx, _ := r.DB.Begin()
	q := `
	UPDATE members SET 
		id=?, name=?, dob=julianday(?), personnummer=?, email=?, mobile=?, 
		address1=?, address2=?, postnummer=?, poststed=?, 
		synagogue_seat=?, membership_fee_tier=?, registered_date=julianday(?), deregistered_date=julianday(?), 
		receive_email=?, receive_mail=?, receive_hatikva=?, archived=?, status=?, 
		family_uuid=?, family_name=? 
	WHERE uuid=?
	`
	_, err := tx.Exec(q,
		member.ID, member.Name, member.DOB, member.Personnummer, member.Email, member.Mobile,
		member.Address.Address1, member.Address.Address2, member.Address.Postnummer, member.Address.Poststed,
		member.Synagogueseat, member.MembershipFeeTier, member.RegisteredDate, member.DeregisteredDate,
		member.ReceiveEmail, member.ReceiveMail, member.ReceiveHatikva, member.Archived, member.Status,
		member.FamilyUUID, member.FamilyName,
		member.UUID,
	)

	if err != nil {
		slog.Error(err.Error(), "uuid", member.UUID, "name", member.Name)
		tx.Rollback()
		return e.ErrUpdatingMember
	}

	// Remove member from all groups
	_, err = tx.Exec(`DELETE FROM members_groups WHERE member_uuid=?`, member.UUID)
	if err != nil {
		slog.Error(err.Error(), "uuid", member.UUID)
		tx.Rollback()
	}

	// Add member to groups
	groupUUIDs = pruneGroupUUIDs(groupUUIDs)
	for _, groupUUID := range groupUUIDs {
		_, err = tx.Exec(`INSERT INTO members_groups(member_uuid, group_uuid) VALUES(?, ?)`, member.UUID, groupUUID)
		if err != nil {
			slog.Error(err.Error(), "uuid", member.UUID, "groupUUD", groupUUID)
			tx.Rollback()
			return e.ErrAddingGroupForMember
		}
	}

	tx.Commit()
	slog.Info("UpdateMember", "uuid", member.UUID, "name", member.Name)
	return nil
}

// TODO: Maybe we don't need to set the empty groupUUIDs to "0", since there is no FK constraint on the member_uuid column in the members_groups table.
func pruneGroupUUIDs(groupUUIDs []string) []string {
	if len(groupUUIDs) == 0 { // Ensure there is at least one groupUUID, i.e. "0"
		groupUUIDs = append(groupUUIDs, "0")
	} else if len(groupUUIDs) > 1 { // Remove "0" if there are other groupUUIDs
		for i, groupUUID := range groupUUIDs {
			if groupUUID == "0" {
				groupUUIDs = append(groupUUIDs[:i], groupUUIDs[i+1:]...)
				break
			}
		}
	}
	return groupUUIDs
}
