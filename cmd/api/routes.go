package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (app *application) HealthCheckHandler(c echo.Context) error {
	data := envelope{
		"status": "available",
		"system_info": envelope{
			"environment": app.config.env,
			"port":        app.config.port,
		},
	}
	return c.JSON(http.StatusOK, data)
}

func (app *application) routes() *echo.Echo {
	e := echo.New()
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(app.rateLimit())

	e.GET("/health", app.HealthCheckHandler)
	e.POST("/user", app.registerUserHandler)
	e.POST("/security/refreshtoken", app.refreshTokenHandler)

	college := e.Group("/college", app.authenticate())
	{
		college.GET("", app.listAllCollegesHandler)
		college.POST("", app.createCollegeHandler)
		college.GET("/:id", app.getCollegeHandler)
		college.PATCH("/:id", app.updateCollegeHandler)
		college.DELETE("/:id", app.deleteCollegeHandler)
	}

	user := e.Group("/user", app.authenticate())
	{
		user.GET("/college/:id", app.listAllUsersInCollegeHandler)
		user.GET("/:id", app.getUserHandler)
		user.PATCH("", app.updateUserHandler)
		user.DELETE("", app.deleteUserHandler)
	}

	return e
}
