package entities

type UserModel struct {
	ID       string
	Email    string
	Password string
}

type UserCreateModel struct {
	Email    string
	Password string
}