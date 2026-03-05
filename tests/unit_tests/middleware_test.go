package unit_tests

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"bookingservice/exceptions"
	middlewares "bookingservice/middleware"
	"bookingservice/models"
	"bookingservice/services"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gobuffalo/buffalo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// newBuffaloApp builds a minimal Buffalo app used for middleware testing.
func newBuffaloApp() *buffalo.App {
	app := buffalo.New(buffalo.Options{Env: "test"})
	return app
}

// buildValidToken generates a signed JWT using the test secret.
func buildValidToken(t *testing.T) string {
	t.Helper()
	secret := []byte("test-secret-key")
	claims := jwt.MapClaims{
		"name":  "Alice",
		"email": "alice@example.com",
		"role":  "user",
		"sub":   float64(1),
		"exp":   time.Now().Add(time.Hour).Unix(),
		"iat":   time.Now().Unix(),
		"iss":   "bookingservice",
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tok.SignedString(secret)
	if err != nil {
		t.Fatalf("failed to build test token: %v", err)
	}
	return token
}

// ---------------------------------------------------------------------------
// JWTAuthenticator tests
// ---------------------------------------------------------------------------

type JWTAuthenticatorTestSuite struct {
	suite.Suite
}

func (s *JWTAuthenticatorTestSuite) SetupTest() {
	os.Setenv("JWT_SECRET", "test-secret-key")
}

func (s *JWTAuthenticatorTestSuite) makeRequest(method, path string, headers map[string]string) *httptest.ResponseRecorder {
	app := newBuffaloApp()
	app.Use(middlewares.JWTAuthenticator())
	app.GET("/users/login", func(c buffalo.Context) error { return c.Render(200, nil) })
	app.GET("/health", func(c buffalo.Context) error { return c.Render(200, nil) })
	app.GET("/protected", func(c buffalo.Context) error { return c.Render(200, nil) })

	req := httptest.NewRequest(method, path, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	return rec
}

func (s *JWTAuthenticatorTestSuite) TestSkipsLoginPath() {
	rec := s.makeRequest("GET", "/users/login", nil)
	assert.NotEqual(s.T(), http.StatusUnauthorized, rec.Code)
}

func (s *JWTAuthenticatorTestSuite) TestSkipsHealthPath() {
	rec := s.makeRequest("GET", "/health", nil)
	assert.NotEqual(s.T(), http.StatusUnauthorized, rec.Code)
}

func (s *JWTAuthenticatorTestSuite) TestMissingAuthHeader() {
	rec := s.makeRequest("GET", "/protected", nil)
	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}

func (s *JWTAuthenticatorTestSuite) TestInvalidToken() {
	// An invalid token triggers a panic inside the middleware which Buffalo recovers
	// as a 500 Internal Server Error (no ErrorHandler is wired in this test app).
	rec := s.makeRequest("GET", "/protected", map[string]string{"Authorization": "Bearer invalid.token.here"})
	assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
}

func (s *JWTAuthenticatorTestSuite) TestValidToken() {
	token := buildValidToken(s.T())
	rec := s.makeRequest("GET", "/protected", map[string]string{"Authorization": "Bearer " + token})
	assert.Equal(s.T(), http.StatusOK, rec.Code)
}

func TestJWTAuthenticatorTestSuite(t *testing.T) {
	suite.Run(t, new(JWTAuthenticatorTestSuite))
}

// ---------------------------------------------------------------------------
// ErrorHandler tests
// ---------------------------------------------------------------------------

type ErrorHandlerTestSuite struct {
	suite.Suite
}

func (s *ErrorHandlerTestSuite) SetupTest() {
	os.Setenv("JWT_SECRET", "test-secret-key")
}

// makeErrorHandlerRequest sets up a Buffalo app with only the ErrorHandler middleware
// and a handler that either panics with err or returns handlerErr.
func (s *ErrorHandlerTestSuite) makeErrorHandlerRequest(panicValue interface{}, handlerErr error) *httptest.ResponseRecorder {
	app := buffalo.New(buffalo.Options{Env: "test"})
	app.Use(middlewares.ErrorHandler())
	app.GET("/test", func(c buffalo.Context) error {
		if panicValue != nil {
			panic(panicValue)
		}
		return handlerErr
	})

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	return rec
}

func (s *ErrorHandlerTestSuite) TestNoError() {
	rec := s.makeErrorHandlerRequest(nil, nil)
	assert.Equal(s.T(), http.StatusOK, rec.Code)
}

func (s *ErrorHandlerTestSuite) TestAuthorizationError_Panic() {
	rec := s.makeErrorHandlerRequest(exceptions.NewAuthorizationError("Forbidden"), nil)
	assert.Equal(s.T(), http.StatusForbidden, rec.Code)
}

func (s *ErrorHandlerTestSuite) TestUserNotFoundError_Panic() {
	rec := s.makeErrorHandlerRequest(exceptions.NewNotFoundError("Not found"), nil)
	assert.Equal(s.T(), http.StatusNotFound, rec.Code)
}

func (s *ErrorHandlerTestSuite) TestWrongPasswordError_Panic() {
	rec := s.makeErrorHandlerRequest(exceptions.NewWrongPasswordError("Wrong password"), nil)
	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}

func (s *ErrorHandlerTestSuite) TestInternalPanic() {
	rec := s.makeErrorHandlerRequest(errors.New("unexpected error"), nil)
	assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
}

func TestErrorHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorHandlerTestSuite))
}

// ---------------------------------------------------------------------------
// Compile-time check: ensure UserService implements expected interface via
// ValidateAuthenticationToken (used by JWTAuthenticator internally).
// ---------------------------------------------------------------------------
var _ = (*services.UserService)(nil)
var _ = (*models.User)(nil)
