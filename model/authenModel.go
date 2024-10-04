package model

type AuthenPayload struct {
	UserId   string `json:"user_id" db:"user_id"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password,omitempty" db:"password"`
	Role     Roles  `json:"role" db:"role"`
}

type AuthenDetial struct {
	UserId  string `json:"user_id" db:"user_id"`
	Name    string `json:"name" db:"name"`
	Surname string `json:"surname" db:"surname"`
	Email   string `json:"email" db:"email"`
	Role    Roles  `json:"role" db:"role"`
}
type AuthenResponse struct {
	Data  AuthenDetial `json:"data"`
	Token string       `json:"token,omitempty"`
}

type AuthenUsecase interface {
	SignIn(user *AuthenPayload) (*AuthenPayload, error)
}

type AuthenRepository interface {
	SignIn(user *AuthenPayload) (*AuthenPayload, error)
}
