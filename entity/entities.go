package entity

type User struct {
	UserName       string
	HashedPassword string
}

func NewUser(username, password string) User {
	return User{
		UserName:       username,
		HashedPassword: password,
	}
}
