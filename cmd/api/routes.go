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
	e.POST("/college", app.createCollegeHandler)
	e.GET("/college/:id", app.getCollegeHandler)
	e.PATCH("/college/:id", app.updateCollegeHandler)
	e.DELETE("/college/:id", app.deleteCollegeHandler)
	return e
}
