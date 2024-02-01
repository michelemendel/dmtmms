package auth

type TokenType string

type UserSession struct {
	Name  string
	Token TokenType
}
