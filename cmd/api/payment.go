package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *application) checkPayementHandler(c echo.Context) error {
	options := map[string]interface{}{}

	body, err := app.razor.Order.All(options, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, envelope{"error in payment": err.Error()})
	}
	return c.JSON(http.StatusOK, body)
}

func (app *application) createPaymentHandler(c echo.Context) error {
	options := map[string]interface{}{
		"amount":          100,
		"currency":        "INR",
		"receipt":         "rcptid_11",
		"payment_capture": 1,
		"notes": map[string]interface{}{
			"name": "Priyan Kishore",
		},
	}

	body, err := app.razor.Order.Create(options, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, envelope{"error in payment": err.Error()})
	}
	return c.JSON(http.StatusOK, body)
}
