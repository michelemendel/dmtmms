package view

import (
	"github.com/michelemendel/dmtmms/entity"
	"github.com/michelemendel/dmtmms/util"
)

// TODO: This should be placed in a view package
// To be presented on the detail section of the member page
// Maybe not the best way to get a list of data in order, when we have to write memberDetails[2].Value to get some value.
func MemberDetailsForPresentation(member entity.Member, groups []entity.Group) entity.MemberDetails {
	if member.UUID == "" {
		return entity.MemberDetails{
			MemberDetails: []entity.MemberDetail{},
			Groups:        []entity.Group{},
		}
	}

	details := []entity.MemberDetail{}
	details = append(details, entity.MemberDetail{Title: "UUID", Value: member.UUID})
	details = append(details, entity.MemberDetail{Title: "FamilyUUID", Value: member.FamilyUUID})
	details = append(details, entity.MemberDetail{Title: "FamilyName", Value: member.FamilyName})
	details = append(details, entity.MemberDetail{Title: "Name", Value: member.Name})
	details = append(details, entity.MemberDetail{Title: "ID", Value: util.Int2String(member.ID)})
	details = append(details, entity.MemberDetail{Title: "Date of Birth", Value: util.Date2String(member.DOB)})
	details = append(details, entity.MemberDetail{Title: "Personnummer", Value: member.Personnummer})
	details = append(details, entity.MemberDetail{Title: "Email", Value: member.Email})
	details = append(details, entity.MemberDetail{Title: "Mobile", Value: member.Mobile})
	details = append(details, entity.MemberDetail{Title: "Address1", Value: member.Address.Address1})
	details = append(details, entity.MemberDetail{Title: "Address2", Value: member.Address.Address2})
	details = append(details, entity.MemberDetail{Title: "Poststed", Value: member.Address.Postnummer + " " + member.Address.Poststed})
	details = append(details, entity.MemberDetail{Title: "Status", Value: string(member.Status)})
	details = append(details, entity.MemberDetail{Title: "Synagogueseat", Value: member.Synagogueseat})
	details = append(details, entity.MemberDetail{Title: "MembershipFeeTier", Value: member.MembershipFeeTier})
	details = append(details, entity.MemberDetail{Title: "RegisteredDate", Value: util.Date2String(member.RegisteredDate)})
	details = append(details, entity.MemberDetail{Title: "DeregisteredDate", Value: util.Date2String(member.DeregisteredDate)})
	details = append(details, entity.MemberDetail{Title: "ReceiveEmail", Value: util.Bool2String(member.ReceiveEmail)})
	details = append(details, entity.MemberDetail{Title: "ReceiveMail", Value: util.Bool2String(member.ReceiveMail)})
	details = append(details, entity.MemberDetail{Title: "ReceiveHatikvah", Value: util.Bool2String(member.ReceiveHatikvah)})
	details = append(details, entity.MemberDetail{Title: "CreatedAt", Value: util.DateTime2String(member.CreatedAt)})
	details = append(details, entity.MemberDetail{Title: "UpdatedAt", Value: util.DateTime2String(member.UpdatedAt)})
	// details = append(details, MemberDetail{"Archived", util.Bool2String(member.Archived)})

	return entity.MemberDetails{
		MemberDetails: details,
		Groups:        groups,
	}
}

func Groups2UUIDsAsStrings(groups []entity.Group) []string {
	groupUUIDs := []string{}
	for _, group := range groups {
		groupUUIDs = append(groupUUIDs, group.UUID)
	}
	return groupUUIDs
}
