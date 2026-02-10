package actions

import (
	"sync"

	middlewares "bookingservice/middleware"
	"bookingservice/services"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/middleware/csrf"
	"github.com/gobuffalo/middleware/i18n"
	"github.com/gobuffalo/middleware/paramlogger"

	"bookingservice/controllers"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")

var (
	app     *buffalo.App
	appOnce sync.Once
	T       *i18n.Translator

	// services
	UserService *services.UserService

	//  controllers
	UserController *controllers.UserController
)

func App() *buffalo.App {
	appOnce.Do(func() {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_bookingservice_session",
		})

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Setup the error handler.
		app.Use(middlewares.JWTAuthenticator())
		app.Use(middlewares.ErrorHandler())

		// health check route
		app.GET("/health", controllers.HealthHandler)

		// user specific routes
		var userRoutes = app.Group("/users")
		userRoutes.GET("/", UserController.GetUsers)
		userRoutes.GET("/{userName}", UserController.GetUser)
		userRoutes.POST("/login", UserController.Login)

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)
	})

	return app
}
