package entity

// DTO
type User struct {
	Name           string
	HashedPassword string
	Role           string
}

func NewUser(name, password, role string) User {
	return User{
		Name:           name,
		HashedPassword: password,
		Role:           role,
	}
}

type Member struct {
	UUID        string
	ID          string
	Name        string
	DateOfBirth string
	Groups      []Group
}

type Group struct {
	UUID    string
	Name    string
	Type    string
	Members []Member
}
