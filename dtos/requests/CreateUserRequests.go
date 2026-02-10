package dtos

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	UserName string `json:"username" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}
