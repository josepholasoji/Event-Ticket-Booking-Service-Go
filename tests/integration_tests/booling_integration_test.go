package controllers_tests

import (
	"bookingservice/actions"
	"bookingservice/controllers"
	"bookingservice/initializations"
	"bookingservice/services"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
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
	initializations.ConnectionToMySQLDB()

	actions.UserService = services.NewUserService()
	actions.UserController = controllers.NewLoginUserController(actions.UserService)

	app := actions.App()
	s.server = httptest.NewServer(app)
	s.baseURL = s.server.URL
}

func (s *BookingControllerTestSuite) TearDownSuite() {
	s.server.Close()
	initializations.CloseMySQLDB()
}

// TestHealthCheck verifies the /health endpoint is reachable and returns 200.
func (s *BookingControllerTestSuite) TestHealthCheck() {
	resp, err := http.Get(s.baseURL + "/health")
	if err != nil {
		s.T().Fatalf("failed to call /health: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.T().Fatalf("expected 200 got %d", resp.StatusCode)
	}
}

// TODO: Add booking-specific integration tests here once booking routes are implemented.

func TestBookingControllerTestSuite(t *testing.T) {
	wd, _ := os.Getwd()
	fmt.Println("WORKDIR:", wd)

	godotenv.Load(".env.integration")
	fmt.Println("testing an environment variable:", os.Getenv("MYSQL_HOST"))
	suite.Run(t, new(BookingControllerTestSuite))
}

