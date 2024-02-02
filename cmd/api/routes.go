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

	e.GET("/", app.HealthCheckHandler)

	college := e.Group("/college")
	{
		college.GET("", app.listAllCollegesHandler)
		college.POST("", app.createCollegeHandler)
		college.GET("/:id", app.getCollegeHandler)
		college.PATCH("/:id", app.updateCollegeHandler)
		college.DELETE("/:id", app.deleteCollegeHandler)
	}

	user := e.Group("/user")
	{
		user.GET("/college/:id", app.listAllUsersInCollegeHandler)
		user.POST("", app.registerUserHandler)
		user.GET("/:id", app.getUserHandler)
		user.PATCH("/:id", app.updateUserHandler)
		user.DELETE("/:id", app.deleteUserHandler)
	}

	return e
}
