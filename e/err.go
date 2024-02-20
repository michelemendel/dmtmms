package e

import "errors"

var (
	ErrSystem                = errors.New("system error")
	ErrCreatingMember        = errors.New("couldn't create member")
	ErrAddingGroupForMember  = errors.New("couldn't add group for member")
	ErrUpdatingMember        = errors.New("couldn't update member")
	ErrUserExists            = errors.New("user name already exists")
	ErrFamilyExists          = errors.New("family name already exists")
	ErrFamilyIsUsedByMembers = errors.New("This family is used by one or more members. Please remove or change the family from the members first.")
	ErrGroupExists           = errors.New("group name already exists")
	ErrGroupIsUsedByMembers  = errors.New("This group is used by one or more members. Please remove or change the group from the members first.")
	ErrInvalidCredentials    = errors.New("invalid credentials")
)
