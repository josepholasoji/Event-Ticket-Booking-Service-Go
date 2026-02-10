package controllers_tests

//import BaseTestSuite

import (
	"bookingservice/actions"
	"bookingservice/controllers"
	"bookingservice/initializations"
	"bookingservice/services"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type BookingControllerTestSuite struct {
	suite.Suite
	baseURL string
	server  *httptest.Server
}

func (s *BookingControllerTestSuite) SetupSuite() {
	// Set up a dummy data for testing if needed
	initializations.ConnectionToMySQLDB()

	actions.UserService = services.NewUserService()
	actions.UserController = controllers.NewLoginUserController(actions.UserService)

	// Start the test server
	app := actions.App()
	s.server = httptest.NewServer(app)
	s.baseURL = s.server.URL + "/users"
}

func (s *BookingControllerTestSuite) TearDownSuite() {
	s.server.Close()
	initializations.CloseMySQLDB()
}

func (s *BookingControllerTestSuite) TestUserLogin_Success() {
	// Ensure that the test user exists in the database before running this test
	var hashedPassword, _ = actions.UserService.HashPassword("testpassword")
	initializations.MySQLDB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", "Test User", "testuser@gmail.com", hashedPassword)

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

	// Clean up the test user after the test
	initializations.MySQLDB.Exec("DELETE FROM users WHERE email = ?", "testuser@gmail.com")
}

func (s *BookingControllerTestSuite) TestUserLogin_Failure_WrongPassword() {
	// Ensure that the test user exists in the database before running this test
	var hashedPassword, _ = actions.UserService.HashPassword("testpassword")
	initializations.MySQLDB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", "Test User", "testuser@gmail.com", hashedPassword)

	// Now perform the login test
	var loginRequest = `{"username": "testuser@gmail.com", "password": "wrongpassword"}`
	resp, err := http.Post(s.baseURL+"/login", "application/json", strings.NewReader(loginRequest))
	if err != nil {
		s.T().Fatalf("failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	// Clean up the test user after the test
	initializations.MySQLDB.Exec("DELETE FROM users WHERE email = ?", "testuser@gmail.com")

	if resp.StatusCode != http.StatusUnauthorized {
		s.T().Fatalf("expected 401 got %d", resp.StatusCode)
	}
}

func (s *BookingControllerTestSuite) TestUserLogin_Failure_UsernameShouldBeEmail() {
	// Ensure that the test user exists in the database before running this test
	var hashedPassword, _ = actions.UserService.HashPassword("testpassword")
	initializations.MySQLDB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", "Test User", "testuser@gmail.com", hashedPassword)

	// Now perform the login test
	var loginRequest = `{"username": "testuser", "password": "wrongpassword"}`
	resp, err := http.Post(s.baseURL+"/login", "application/json", strings.NewReader(loginRequest))
	if err != nil {
		s.T().Fatalf("failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	// Clean up the test user after the test
	initializations.MySQLDB.Exec("DELETE FROM users WHERE email = ?", "testuser@gmail.com")

	if resp.StatusCode != http.StatusBadRequest {
		s.T().Fatalf("expected 400 got %d", resp.StatusCode)
	}
}

func (s *BookingControllerTestSuite) TestUserLogin_Failure_UserNotFound() {
	// Ensure that the test user exists in the database before running this test
	var hashedPassword, _ = actions.UserService.HashPassword("testpassword")
	initializations.MySQLDB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", "Test User", "testuser@gmail.com", hashedPassword)

	// Now perform the login test
	var loginRequest = `{"username": "testuserX@gmail.com", "password": "testpassword"}`
	resp, err := http.Post(s.baseURL+"/login", "application/json", strings.NewReader(loginRequest))
	if err != nil {
		s.T().Fatalf("failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	// Clean up the test user after the test
	initializations.MySQLDB.Exec("DELETE FROM users WHERE email = ?", "testuser@gmail.com")

	if resp.StatusCode != http.StatusNotFound {
		s.T().Fatalf("expected 404 got %d", resp.StatusCode)
	}
}

func TestBookingControllerTestSuite(t *testing.T) {
	wd, _ := os.Getwd()
	fmt.Println("WORKDIR:", wd)

	godotenv.Load(".env.integration")
	fmt.Println("testing an environment variable:", os.Getenv("MYSQL_HOST"))
	suite.Run(t, new(BookingControllerTestSuite))
}
