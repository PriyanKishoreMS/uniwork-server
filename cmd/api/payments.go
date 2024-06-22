package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/priyankishorems/uniwork-server/internal/data"
)

func (app *application) checkPayementHandler(c echo.Context) error {
	options := map[string]interface{}{}

	body, err := app.razor.Order.All(options, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, envelope{"error in payment": err.Error()})
	}
	return c.JSON(http.StatusOK, body)
}

func (app *application) createOrderHandler(c echo.Context) error {
	var input struct {
		Amount   int64 `json:"amount" validate:"required"`
		TaskID   int64 `json:"task_id" validate:"required"`
		WorkerID int64 `json:"worker_id" validate:"required"`
	}

	requester := app.contextGetUser(c)

	if err := app.readJSON(c, &input); err != nil {
		app.BadRequest(c, err)
		return err
	}

	err := app.validate.Struct(input)
	if err != nil {
		app.ValidationError(c, err)
		return err
	}

	storedTaskOwnerId, storedTaskPrice, _, err := app.models.Tasks.GetTaskForVerification(input.TaskID)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	if storedTaskOwnerId != requester.ID {
		app.UserUnAuthorizedResponse(c)
		return ErrUserUnauthorized
	}

	if storedTaskPrice != int64(input.Amount) {
		app.UserUnAuthorizedResponse(c)
		return ErrUserUnauthorized
	}

	options := map[string]interface{}{
		"amount":          input.Amount,
		"currency":        "INR",
		"payment_capture": 1,
		"notes": map[string]interface{}{
			"name":  requester.Name,
			"email": requester.Email,
		},
	}

	body, err := app.razor.Order.Create(options, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, envelope{"error in payment": err.Error()})
	}

	orderData := data.OrderData{
		OrderID:     body["id"].(string),
		Amount:      input.Amount,
		TaskID:      input.TaskID,
		TaskOwnerID: requester.ID,
		WorkerID:    input.WorkerID,
	}

	err = app.models.Payments.CreateOrder(&orderData)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	return c.JSON(http.StatusOK, envelope{"orderId": orderData.OrderID, "amount": orderData.Amount})
}
