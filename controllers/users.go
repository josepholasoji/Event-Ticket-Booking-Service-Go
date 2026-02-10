package controllers

import (
	dtos "bookingservice/dtos/requests"
	"bookingservice/models"
	"bookingservice/services"
	"bookingservice/utils"
	"encoding/json"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gobuffalo/buffalo"
)

type UserController struct {
	services *services.UserService
}

func (u UserController) GetUsers(c buffalo.Context) error {
	users, err := u.services.GetAllUsers()
	if err != nil {
		return utils.DoResponse(c, 400, 1, "Failed to fetch users: "+err.Error(), "")
	}
	str, _ := json.Marshal(users)
	return utils.DoResponse(c, 200, 0, "Users fetched successfully", string(str))
}

func (u UserController) GetUser(c buffalo.Context) error {
	// Get the user info from the context set by the JWT middleware
	var userInfo = c.Value("user_info")
	var user = userInfo.(*models.User)
	if user.ID == 0 {
		return utils.DoResponse(c, 401, 1, "Unauthorized", "")
	}

	var userName = c.Param("userName")
	if user.UserName != userName && strings.ToLower(user.Role) != "admin" {
		return utils.DoResponse(c, 403, 1, "Forbidden: You can only access your own user details", "")
	}

	userFromDB, err := u.services.GetUserByUserName(userName)
	if err != nil {
		return utils.DoResponse(c, 400, 1, "Failed to fetch user: "+err.Error(), "")
	}

	str, _ := json.Marshal(userFromDB)
	return utils.DoResponse(c, 200, 0, "User fetched successfully", string(str))
}

func (u UserController) Login(c buffalo.Context) error {
	var loginUserRequest dtos.LoginUserRequest
	if err := c.Bind(&loginUserRequest); err != nil {
		return utils.DoResponse(c, 400, 1, "Invalid request body", "")
	}

	validator := validator.New()
	if err := validator.Struct(loginUserRequest); err != nil {
		return utils.DoResponse(c, 400, 1, "Validation failed: "+err.Error(), "")
	}

	user, err := u.services.LoginUser(loginUserRequest)
	if err != nil {
		return utils.DoResponse(c, 400, 1, "Login failed: "+err.Error(), "")
	}

	token, err := u.services.GenerateToken(user)
	if err != nil {
		return utils.DoResponse(c, 500, 1, "Token generation failed: "+err.Error(), "")
	}
	return utils.DoResponse(c, 200, 0, "Login successful", token)
}

func NewLoginUserController(usr *services.UserService) *UserController {
	return &UserController{
		services: usr,
	}
}
