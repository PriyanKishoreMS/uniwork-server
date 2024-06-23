package main

import (
	"net/http"
	"os"

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
		Amount        int64 `json:"amount" validate:"required"`
		TaskID        int64 `json:"task_id" validate:"required"`
		TaskRequestID int64 `json:"task_request_id" validate:"required"`
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

	storedTaskOwnerId, storedTaskPrice, storedTaskId, err := app.models.TaskRequests.OrderCreationCheck(input.TaskRequestID)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	// check if the requester is the owner of the task
	// check if the price of the task is same as the input amount
	// check if the task request belongs to the task
	if storedTaskOwnerId != requester.ID || storedTaskPrice != int64(input.Amount) || storedTaskId != input.TaskID {
		app.UserUnAuthorizedResponse(c)
		return ErrUserUnauthorized
	}

	// if the order already exists return the order id
	existingRazorPayOrderId := app.models.Payments.CheckExistingOrder(input.TaskRequestID)
	if existingRazorPayOrderId != nil {
		return c.JSON(http.StatusOK, envelope{"orderId": *existingRazorPayOrderId})
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
		OrderID:       body["id"].(string),
		Amount:        input.Amount,
		TaskOwnerID:   requester.ID,
		TaskRequestID: input.TaskRequestID,
	}

	err = app.models.Payments.CreateOrder(&orderData)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	return c.JSON(http.StatusOK, envelope{"orderId": orderData.OrderID})
}

// TODO Setup webhook to handle payment status
func (app *application) PaymentStatusHandler(c echo.Context) error {

	var input struct {
		OrderID          string `json:"razorpay_order_id" validate:"required"`
		PaymentID        string `json:"razorpay_payment_id" validate:"required"`
		PaymentSignature string `json:"razorpay_signature" validate:"required"`
	}

	if err := app.readJSON(c, &input); err != nil {
		app.BadRequest(c, err)
		return err
	}
	secret := os.Getenv("RAZORPAY_SECRET")

	generatedSignature := generateSignature(input.OrderID, input.PaymentID, secret)
	if generatedSignature != input.PaymentSignature {
		return c.JSON(http.StatusBadRequest, envelope{"error": "Invalid Payment"})
	}

	// TODO approve the task request
	// select tr.task_id, tr.task_worker_id, t.version from payment_details p join task_requests tr on p.task_request_id=tr.id join tasks t on t.id=tr.task_id where p.razorpay_order_id='order_OQHkSZwDQmzGhD';

	return c.JSON(http.StatusOK, envelope{"status": "Payment Successful"})
}
