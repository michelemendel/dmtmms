package repo

import (
	"log/slog"

	// "github.com/michelemendel/dmtmms/e"
	"github.com/michelemendel/dmtmms/e"
	"github.com/michelemendel/dmtmms/entity"
)

func (r *Repo) CreateMember(member entity.Member, groupUUIDs []string) error {
	tx, _ := r.DB.Begin()
	_, err := tx.Exec(`
	INSERT INTO members(
		uuid, 
		id, 
		name, 
		dob, 
		personnummer, 
		email, 
		mobile, 
		address1, 
		address2, 
		postnummer, 
		poststed, 
		synagogue_seat,
		membership_fee_tier,
		registered_date,
		deregistered_date,
		receive_email,
		receive_mail,
		receive_hatikva,
		archived,
		status,
		family_uuid,
		family_group
	) VALUES(?, ?, ?, julianday(?), ?, ?, ?, ?, ?, ?, ?, ?, ?, julianday(?), julianday(?), ?, ?, ?, ?, ?, ?, ?)
	`,
		member.UUID,
		member.ID,
		member.Name,
		member.DOB,
		member.Personnummer,
		member.Email,
		member.Mobile,
		member.Address.Address1,
		member.Address.Address2,
		member.Address.Postnummer,
		member.Address.Poststed,
		member.Synagogueseat,
		member.MembershipFeeTier,
		member.RegisteredDate,
		member.DeregisteredDate,
		member.ReceiveEmail,
		member.ReceiveMail,
		member.ReceiveHatikva,
		member.Archived,
		member.Status,
		member.FamilyUUID,
		member.FamilyGroup,
	)
	if err != nil {
		slog.Error(err.Error(), "uuid", member.UUID, "name", member.Name)
		tx.Rollback()
		return e.ErrCreatingMember
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
	_, err := r.DB.Exec("DELETE FROM members WHERE uuid=?", memberUUID)
	if err != nil {
		slog.Error(err.Error(), "uuid", memberUUID)
		return err
	}
	slog.Info("DeleteMember", "uuid", memberUUID)
	return nil
}

func (r *Repo) UpdateMember(member entity.Member) error {
	slog.Info("UpdateGroup", "uuid", member.UUID, "name", member.Name)
	return nil
}
