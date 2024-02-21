package repo

import (
	"fmt"
	"log/slog"

	// "github.com/michelemendel/dmtmms/e"
	"github.com/michelemendel/dmtmms/e"
	"github.com/michelemendel/dmtmms/entity"
)

func (r *Repo) CreateMember(member entity.Member, groupUUIDs []string) error {
	tx, _ := r.DB.Begin()
	familyUUID := "0"
	familyName := ""
	if member.FamilyUUID != "" {
		familyUUID = member.FamilyUUID
		familyName, _ = r.GetFamilyNameByUUID(member.FamilyUUID)
	}

	// fmt.Println("[R]: familyUUID", familyUUID, "familyName", familyName)

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
	gUUIDs := []string{"0"}
	if len(groupUUIDs) > 0 {
		gUUIDs = groupUUIDs
	}

	fmt.Println("[R]: groupUUIDs", groupUUIDs, "len", len(groupUUIDs), "gUUIDs", gUUIDs)

	for _, groupUUID := range gUUIDs {
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
	fmt.Println("[R]: ArchiveMember", "memberUUID", memberUUID)
	_, err := r.DB.Exec("UPDATE members SET archived=true WHERE uuid=?", memberUUID)
	if err != nil {
		slog.Error(err.Error(), "uuid", memberUUID)
		return err
	}
	slog.Info("ArchivedMember", "uuid", memberUUID)
	return nil
}

func (r *Repo) DeleteMember(memberUUID string) error {
	fmt.Println("[R]: DeleteMember", "memberUUID", memberUUID)
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
func (r *Repo) UpdateMember(member entity.Member, groupUUIDs []string) error {
	tx, _ := r.DB.Begin()
	_, err := tx.Exec(`
	UPDATE members SET 
		id=?, name=?, dob=julianday(?), personnummer=?, email=?, mobile=?, 
		address1=?, address2=?, postnummer=?, poststed=?, 
		synagogue_seat=?, membership_fee_tier=?, registered_date=julianday(?), deregistered_date=julianday(?), 
		receive_email=?, receive_mail=?, receive_hatikva=?, archived=?, status=?, 
		family_uuid=?, family_name=? 
	WHERE uuid=?
	`,
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
