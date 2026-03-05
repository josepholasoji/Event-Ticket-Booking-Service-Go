package unit_tests

import (
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"

	dtos "bookingservice/dtos/requests"
	"bookingservice/models"
	"bookingservice/services"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// mockUserRepository is a test-double for UserRepositoryInterface.
type mockUserRepository struct {
	users    []models.User
	saveErr  error
	findErr  error
	getAllErr error
}

func (m *mockUserRepository) GetAllUsers() ([]models.User, error) {
	return m.users, m.getAllErr
}

func (m *mockUserRepository) FindUserByName(username string) (*models.User, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	for i := range m.users {
		if m.users[i].UserName == username {
			return &m.users[i], nil
		}
	}
	return nil, sql.ErrNoRows
}

func (m *mockUserRepository) SaveUser(name, username, password string) error {
	return m.saveErr
}

// ---------------------------------------------------------------------------

type UserServiceTestSuite struct {
	suite.Suite
	svc  *services.UserService
	repo *mockUserRepository
}

func (s *UserServiceTestSuite) SetupTest() {
	os.Setenv("JWT_SECRET", "test-secret-key")
	s.repo = &mockUserRepository{}
	s.svc = services.NewUserServiceWithRepository(s.repo)
}

// ---------------------------------------------------------------------------
// HashPassword / CheckPasswordHash
// ---------------------------------------------------------------------------

func (s *UserServiceTestSuite) TestHashPassword() {
	hash, err := s.svc.HashPassword("mypassword")
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), hash)
	assert.NotEqual(s.T(), "mypassword", hash)
}

func (s *UserServiceTestSuite) TestCheckPasswordHash_Valid() {
	hash, _ := s.svc.HashPassword("mypassword")
	assert.True(s.T(), s.svc.CheckPasswordHash("mypassword", hash))
}

func (s *UserServiceTestSuite) TestCheckPasswordHash_Invalid() {
	hash, _ := s.svc.HashPassword("mypassword")
	assert.False(s.T(), s.svc.CheckPasswordHash("wrongpassword", hash))
}

// ---------------------------------------------------------------------------
// GenerateToken / ValidateAuthenticationToken
// ---------------------------------------------------------------------------

func (s *UserServiceTestSuite) TestGenerateToken() {
	user := &models.User{ID: 1, Name: "Alice", UserName: "alice@example.com", Role: "user"}
	token, err := s.svc.GenerateToken(user)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), token)
}

func (s *UserServiceTestSuite) TestValidateAuthenticationToken_Valid() {
	user := &models.User{ID: 1, Name: "Alice", UserName: "alice@example.com", Role: "user"}
	token, _ := s.svc.GenerateToken(user)

	result, err := s.svc.ValidateAuthenticationToken(token)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user.UserName, result.UserName)
	assert.Equal(s.T(), user.Name, result.Name)
	assert.Equal(s.T(), user.Role, result.Role)
}

func (s *UserServiceTestSuite) TestValidateAuthenticationToken_Invalid() {
	_, err := s.svc.ValidateAuthenticationToken("not.a.valid.token")
	assert.Error(s.T(), err)
}

func (s *UserServiceTestSuite) TestValidateAuthenticationToken_Expired() {
	// Manually build a token that expired in the past.
	secret := []byte("test-secret-key")
	claims := jwt.MapClaims{
		"name":  "Alice",
		"email": "alice@example.com",
		"role":  "user",
		"sub":   float64(1),
		"exp":   time.Now().Add(-time.Hour).Unix(),
		"iat":   time.Now().Add(-2 * time.Hour).Unix(),
		"iss":   "bookingservice",
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Use WithoutVerification to sign with the correct secret but expired exp
	tokenString, _ := tok.SignedString(secret)

	_, err := s.svc.ValidateAuthenticationToken(tokenString)
	assert.Error(s.T(), err)
}

func (s *UserServiceTestSuite) TestValidateAuthenticationTokenNoexpiryCheck_Valid() {
	// Generate a normal non-expired token and verify the noexpiry variant also accepts it.
	user := &models.User{ID: 1, Name: "Alice", UserName: "alice@example.com", Role: "user"}
	token, _ := s.svc.GenerateToken(user)

	result, err := s.svc.ValidateAuthenticationTokenNoexpiryCheck(token)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "alice@example.com", result.UserName)
}

// ---------------------------------------------------------------------------
// CreateUser
// ---------------------------------------------------------------------------

func (s *UserServiceTestSuite) TestCreateUser_Success() {
	s.repo.saveErr = nil
	req := dtos.CreateUserRequest{Name: "Bob", UserName: "bob@example.com", Password: "secret123"}
	err := s.svc.CreateUser(req)
	assert.NoError(s.T(), err)
}

func (s *UserServiceTestSuite) TestCreateUser_RepositoryError() {
	s.repo.saveErr = errors.New("db error")
	req := dtos.CreateUserRequest{Name: "Bob", UserName: "bob@example.com", Password: "secret123"}
	err := s.svc.CreateUser(req)
	assert.Error(s.T(), err)
}

// ---------------------------------------------------------------------------
// LoginUser
// ---------------------------------------------------------------------------

func (s *UserServiceTestSuite) TestLoginUser_Success() {
	hash, _ := s.svc.HashPassword("secret123")
	s.repo.users = []models.User{{ID: 1, Name: "Bob", UserName: "bob@example.com", Password: hash, Role: "user"}}

	req := dtos.LoginUserRequest{Username: "bob@example.com", Password: "secret123"}
	user, err := s.svc.LoginUser(req)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "bob@example.com", user.UserName)
}

func (s *UserServiceTestSuite) TestLoginUser_UserNotFound() {
	s.repo.users = []models.User{} // empty — FindUserByName returns sql.ErrNoRows

	assert.Panics(s.T(), func() {
		req := dtos.LoginUserRequest{Username: "nobody@example.com", Password: "secret123"}
		_, _ = s.svc.LoginUser(req)
	})
}

func (s *UserServiceTestSuite) TestLoginUser_WrongPassword() {
	hash, _ := s.svc.HashPassword("secret123")
	s.repo.users = []models.User{{ID: 1, Name: "Bob", UserName: "bob@example.com", Password: hash, Role: "user"}}

	assert.Panics(s.T(), func() {
		req := dtos.LoginUserRequest{Username: "bob@example.com", Password: "wrongpassword"}
		_, _ = s.svc.LoginUser(req)
	})
}

// ---------------------------------------------------------------------------
// GetAllUsers
// ---------------------------------------------------------------------------

func (s *UserServiceTestSuite) TestGetAllUsers_Success() {
	s.repo.users = []models.User{
		{ID: 1, Name: "Alice", UserName: "alice@example.com", Role: "admin"},
		{ID: 2, Name: "Bob", UserName: "bob@example.com", Role: "user"},
	}
	users, err := s.svc.GetAllUsers()
	assert.NoError(s.T(), err)
	assert.Len(s.T(), users, 2)
}

func (s *UserServiceTestSuite) TestGetAllUsers_Error() {
	s.repo.getAllErr = errors.New("db error")
	users, err := s.svc.GetAllUsers()
	assert.Error(s.T(), err)
	assert.Nil(s.T(), users)
}

// ---------------------------------------------------------------------------
// GetUserByUserName
// ---------------------------------------------------------------------------

func (s *UserServiceTestSuite) TestGetUserByUserName_Success() {
	s.repo.users = []models.User{{ID: 1, Name: "Alice", UserName: "alice@example.com", Password: "hash", Role: "user"}}

	user, err := s.svc.GetUserByUserName("alice@example.com")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "alice@example.com", user.UserName)
	assert.Empty(s.T(), user.Password, "password hash must be redacted in response")
}

func (s *UserServiceTestSuite) TestGetUserByUserName_NotFound() {
	s.repo.users = []models.User{} // FindUserByName returns sql.ErrNoRows

	user, err := s.svc.GetUserByUserName("ghost@example.com")
	assert.Error(s.T(), err)
	assert.Nil(s.T(), user)
}

// ---------------------------------------------------------------------------

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
