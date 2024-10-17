package model

type AuthenPayload struct {
	UserId   string `json:"user_id"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Role     Roles  `json:"role"`
}

type AuthenDetail struct {
	UserId  string `json:"user_id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email"`
	Role    Roles  `json:"role"`
	Phone   string `json:"phone"`
}
type AuthenResponse struct {
	Data  AuthenDetail `json:"data"`
	Token string       `json:"token,omitempty"`
}

type AuthenUsecase interface {
	SignIn(user *AuthenPayload) (*AuthenPayload, error)
	Me(userId string) (*AuthenResponse, error)
}

type AuthenRepository interface {
	SignIn(user *AuthenPayload) (*AuthenPayload, error)
	Me(userId string) (*AuthenResponse, error)
}
