package handlers

type UserDTO struct {
	Id       string `json:"id"`
	Email    string `json:"email" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Agent    bool   `json:"agent" validate:"required"`
	Token    string `json:"token"`
}

type SignInDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
