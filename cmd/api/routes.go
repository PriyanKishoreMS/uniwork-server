package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type wrap map[string]interface{}

func (app *application) HealthCheckHandler(c echo.Context) error {
	data := wrap{
		"status": "available",
		"system_info": wrap{
			"environment": app.config.env,
			"port":        app.config.port,
		},
	}
	return c.JSON(http.StatusOK, data)
}

func (app *application) routes(e *echo.Echo) {
	e.GET("/", app.HealthCheckHandler)
}
