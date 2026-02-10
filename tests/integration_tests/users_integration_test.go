package controllers_tests

//import BaseTestSuite

import (
	"bookingservice/actions"
	"bookingservice/controllers"
	"bookingservice/dtos/responses"
	"bookingservice/initializations"
	"bookingservice/models"
	"bookingservice/services"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type UsersControllerTestSuite struct {
	suite.Suite
	baseURL string
	server  *httptest.Server
}

func (s *UsersControllerTestSuite) SetupSuite() {
	// Set up a dummy data for testing if needed
	initializations.ConnectionToMySQLDB()

	actions.UserService = services.NewUserService()
	actions.UserController = controllers.NewLoginUserController(actions.UserService)

	// Start the test server
	app := actions.App()
	s.server = httptest.NewServer(app)
	s.baseURL = s.server.URL + "/users"
}

func (s *UsersControllerTestSuite) TearDownSuite() {
	s.server.Close()
	initializations.CloseMySQLDB()
}

func (s *UsersControllerTestSuite) TestUserLogin_Success() {
	// Ensure that the test user exists in the database before running this test
	var hashedPassword, _ = actions.UserService.HashPassword("testpassword")
	initializations.MySQLDB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", "Test User", "testuser@gmail.com", hashedPassword)
	defer func() {
		initializations.MySQLDB.Exec("DELETE FROM users WHERE email = ?", "testuser@gmail.com")
	}()

	// Now perform the login test
	var loginRequest = `{"username": "testuser@gmail.com", "password": "testpassword"}`
	resp, err := http.Post(s.baseURL+"/login", "application/json", strings.NewReader(loginRequest))
	if err != nil {
		s.T().Fatalf("failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.T().Fatalf("expected 200 got %d", resp.StatusCode)
	}
}

func (s *UsersControllerTestSuite) TestUserGetAllUsers_Success() {
	// Ensure that the test user exists in the database before running this test
	var hashedPassword, _ = actions.UserService.HashPassword("testpassword")
	initializations.MySQLDB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", "Test User", "testuser@gmail.com", hashedPassword)

	defer func() {
		initializations.MySQLDB.Exec("DELETE FROM users WHERE email = ?", "testuser@gmail.com")
	}()

	// Now perform the login test
	var loginRequest = `{"username": "testuser@gmail.com", "password": "testpassword"}`
	resp, err := http.Post(s.baseURL+"/login", "application/json", strings.NewReader(loginRequest))
	if err != nil {
		s.T().Fatalf("failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.T().Fatalf("expected 200 got %d", resp.StatusCode)
	}

	// Parse the response body to extract the token
	var responseBody = responses.Response{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		s.T().Fatalf("failed to decode response body: %v", err)
	}

	// Extract the user data from the returned autoh token response data
	userSerice := services.NewUserService()
	_, err = userSerice.ValidateAuthenticationToken(responseBody.Data)
	if err != nil {
		s.T().Fatalf("failed to validate authentication token: %v", err)
	}

	// call the get user endpoint, setting the Authorization header with an expired token
	var requestToken = responseBody.Data
	req, err := http.NewRequest("GET", s.baseURL, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+requestToken)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.T().Fatalf("expected 200 got %d", resp.StatusCode)
	}
}

func (s *UsersControllerTestSuite) TestUserGetAllUsers_FailureExpiredToken() {
	// call the get user endpoint, setting the Authorization header with an expired token
	var expiredToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3R1c2VyQGdtYWlsLmNvbSIsImV4cCI6MTc3MDc3NjI2MiwiaWF0IjoxNzcwNjg5ODYyLCJpc3MiOiJib29raW5nc2VydmljZSIsIm5hbWUiOiJUZXN0IFVzZXIiLCJyb2xlIjoidXNlciIsInN1YiI6MTE2fQ.sQgx0sTPxr1bgzXvRLpp_v1fvdPWlk3busdpFh_BAUs"
	req, err := http.NewRequest("GET", s.baseURL, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+expiredToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		s.T().Fatalf("expected 401 got %d", resp.StatusCode)
	}
}

func (s *UsersControllerTestSuite) TestUserLoginAndGetAUser_Success() {
	// Ensure that the test user exists in the database before running this test
	var hashedPassword, _ = actions.UserService.HashPassword("testpassword")
	initializations.MySQLDB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", "Test User", "testuser@gmail.com", hashedPassword)
	defer func() {
		initializations.MySQLDB.Exec("DELETE FROM users WHERE email = ?", "testuser@gmail.com")
	}()

	// Now perform the login test
	var loginRequest = `{"username": "testuser@gmail.com", "password": "testpassword"}`
	resp, err := http.Post(s.baseURL+"/login", "application/json", strings.NewReader(loginRequest))
	if err != nil {
		s.T().Fatalf("failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.T().Fatalf("expected 200 got %d", resp.StatusCode)
	}

	// Parse the response body to extract the token
	var responseBody = responses.Response{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		s.T().Fatalf("failed to decode response body: %v", err)
	}

	userService := services.NewUserService()
	usr, err := userService.ValidateAuthenticationToken(responseBody.Data)
	if err != nil {
		s.T().Fatalf("failed to validate authentication token: %v", err)
	}

	// call the get user endpoint, setting the Authorization header with an expired token
	var expiredToken = responseBody.Data
	req, err := http.NewRequest("GET", s.baseURL+"/"+usr.UserName, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+expiredToken)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.T().Fatalf("expected 200 got %d", resp.StatusCode)
	}

	// Parse the response body to extract the user data
	var getUserResponseBody = responses.Response{}
	if err := json.NewDecoder(resp.Body).Decode(&getUserResponseBody); err != nil {
		s.T().Fatalf("failed to decode response body: %v", err)
	}

	var getUserResponseData = models.User{}
	if err := json.Unmarshal([]byte(getUserResponseBody.Data), &getUserResponseData); err != nil {
		s.T().Fatalf("failed to unmarshal user data: %v", err)
	}

	// Assert that the user data in the response matches the expected data
	if getUserResponseData.UserName != usr.UserName {
		s.T().Fatalf("expected username %s got %s", usr.UserName, getUserResponseData.UserName)
	}

	if getUserResponseData.Name != usr.Name {
		s.T().Fatalf("expected name %s got %s", usr.Name, getUserResponseData.Name)
	}

	if getUserResponseData.Role != usr.Role {
		s.T().Fatalf("expected role %s got %s", usr.Role, getUserResponseData.Role)
	}
}

func (s *UsersControllerTestSuite) TestUserLogin_Failure_WrongPassword() {
	// Ensure that the test user exists in the database before running this test
	var hashedPassword, _ = actions.UserService.HashPassword("testpassword")
	initializations.MySQLDB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", "Test User", "testuser@gmail.com", hashedPassword)
	defer func() {
		initializations.MySQLDB.Exec("DELETE FROM users WHERE email = ?", "testuser@gmail.com")
	}()

	// Now perform the login test
	var loginRequest = `{"username": "testuser@gmail.com", "password": "wrongpassword"}`
	resp, err := http.Post(s.baseURL+"/login", "application/json", strings.NewReader(loginRequest))
	if err != nil {
		s.T().Fatalf("failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		s.T().Fatalf("expected 401 got %d", resp.StatusCode)
	}
}

func (s *UsersControllerTestSuite) TestUserLogin_Failure_UsernameShouldBeEmail() {
	// Ensure that the test user exists in the database before running this test
	var hashedPassword, _ = actions.UserService.HashPassword("testpassword")
	initializations.MySQLDB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", "Test User", "testuser@gmail.com", hashedPassword)
	defer func() {
		initializations.MySQLDB.Exec("DELETE FROM users WHERE email = ?", "testuser@gmail.com")
	}()

	// Now perform the login test
	var loginRequest = `{"username": "testuser", "password": "wrongpassword"}`
	resp, err := http.Post(s.baseURL+"/login", "application/json", strings.NewReader(loginRequest))
	if err != nil {
		s.T().Fatalf("failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		s.T().Fatalf("expected 400 got %d", resp.StatusCode)
	}
}

func (s *UsersControllerTestSuite) TestUserLogin_Failure_UserNotFound() {
	// Ensure that the test user exists in the database before running this test
	var hashedPassword, _ = actions.UserService.HashPassword("testpassword")
	initializations.MySQLDB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", "Test User", "testuser@gmail.com", hashedPassword)
	defer func() {
		initializations.MySQLDB.Exec("DELETE FROM users WHERE email = ?", "testuser@gmail.com")
	}()

	// Now perform the login test
	var loginRequest = `{"username": "testuserX@gmail.com", "password": "testpassword"}`
	resp, err := http.Post(s.baseURL+"/login", "application/json", strings.NewReader(loginRequest))
	if err != nil {
		s.T().Fatalf("failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		s.T().Fatalf("expected 404 got %d", resp.StatusCode)
	}
}

func TestUsersControllerTestSuite(t *testing.T) {
	wd, _ := os.Getwd()
	fmt.Println("WORKDIR:", wd)

	godotenv.Load(".env.integration")
	fmt.Println("testing an environment variable:", os.Getenv("MYSQL_HOST"))
	suite.Run(t, new(UsersControllerTestSuite))
}
