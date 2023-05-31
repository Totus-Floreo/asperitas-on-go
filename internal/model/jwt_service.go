package model

type IJWTService interface {
	GenerateToken(*User) (string, error)
	VerifyToken(string) (*Author, error)
}
