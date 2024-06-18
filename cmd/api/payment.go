package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *application) checkPayementHander(c echo.Context) error {
	options := map[string]interface{}{
		"count": 1,
	}

	body, err := app.razor.Order.All(options, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, envelope{"error in payment": err.Error()})
	}
	return c.JSON(http.StatusOK, body)
}
