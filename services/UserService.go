package services

import (
	"database/sql"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	dtos "bookingservice/dtos/requests"
	"bookingservice/exceptions"
	"bookingservice/models"
	"bookingservice/repositories"

	"github.com/golang-jwt/jwt/v5"
)

type UserService struct {
	userRepository *repositories.UserRepository
}

func (s *UserService) CreateUser(request dtos.CreateUserRequest) error {
	hashedPassword, err := s.HashPassword(request.Password)
	if err != nil {
		return err
	}
	return s.userRepository.SaveUser(request.Name, request.UserName, hashedPassword)
}

func (s *UserService) LoginUser(request dtos.LoginUserRequest) (*models.User, error) {
	userRecord, err := s.userRepository.FindUserByName(request.Username)
	if err == sql.ErrNoRows {
		panic(exceptions.NewNotFoundError("User not found"))
	}

	if err != nil {
		return nil, err
	}

	if s.CheckPasswordHash(request.Password, userRecord.Password) {
		if err != nil {
			return nil, err
		}
		return userRecord, nil
	}

	panic(exceptions.NewWrongPasswordError("Invalid password"))
}

func (s *UserService) GenerateToken(user *models.User) (string, error) {
	var jwtSecret = []byte(os.Getenv("JWT_SECRET")) // move to env in real apps

	claims := jwt.MapClaims{
		"name":  user.Name,
		"email": user.UserName,
		"role":  user.Role,
		"sub":   user.ID,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"iat":   time.Now().Unix(),
		"iss":   "bookingservice",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (s *UserService) ValidateAuthenticationToken(tokenString string) (*models.User, error) {
	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, exceptions.NewAuthorizationError("Invalid token")
	}
	// Check for expiration
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < jwt.NewNumericDate(time.Now()).Unix() {
				return nil, exceptions.NewAuthorizationError("Token has expired")
			}
		}
	}

	// contruct the user info from the token claims and return it
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		var user = models.User{
			Name:     claims["name"].(string),
			UserName: claims["email"].(string),
			Role:     claims["role"].(string),
			ID:       uint(claims["sub"].(float64)),
		}
		return &user, nil
	}

	return nil, exceptions.NewAuthorizationError("Invalid token claims")
}

func (s *UserService) ValidateAuthenticationTokenNoexpiryCheck(tokenString string) (*models.User, error) {
	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, exceptions.NewAuthorizationError("Invalid token")
	}

	// contruct the user info from the token claims and return it
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		var user = models.User{
			Name:     claims["name"].(string),
			UserName: claims["email"].(string),
			Role:     claims["role"].(string),
			ID:       uint(claims["sub"].(float64)),
		}
		return &user, nil
	}

	return nil, exceptions.NewAuthorizationError("Invalid token claims")
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.userRepository.GetAllUsers()
}

func (s *UserService) GetUserByUserName(userName string) (*models.User, error) {
	user, err := s.userRepository.FindUserByName(userName)
	if err == sql.ErrNoRows {
		return nil, exceptions.NewNotFoundError("User not found")
	}
	user.Password = "" // do not return the password hash
	return user, nil
}

func (s *UserService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	return string(bytes), err
}

func (s *UserService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func NewUserService() *UserService {
	return &UserService{
		userRepository: repositories.NewUserRepository(),
	}
}
