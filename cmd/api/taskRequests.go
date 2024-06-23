package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *application) addNewTaskRequestHandler(c echo.Context) error {
	taskWorkerId, err1 := app.readIntParam(c, "userid")
	taskId, err2 := app.readIntParam(c, "taskid")
	if err1 != nil || err2 != nil {
		app.BadRequest(c, fmt.Errorf("error reading form data: %w", err1))
		return fmt.Errorf("error reading form data: %w", err1)
	}

	fmt.Println(taskWorkerId, taskId, "taskWorkerId, taskId")

	// ! commented for ease of testing, uncomment before deployment
	// requestedUser := app.contextGetUser(c)

	// if requestedUser.ID != taskWorkerId {
	// 	app.CustomErrorResponse(c, envelope{"unauthorized": "You are not authorized to create this task request"}, http.StatusUnauthorized, ErrUserUnauthorized)
	// 	return ErrUserUnauthorized
	// }

	approved, err := app.models.TaskRequests.CheckTaskRequestStatus(taskId)
	if err != nil {
		app.InternalServerError(c, fmt.Errorf("could not check task request status: %w", err))
		return err
	}

	if approved {
		app.CustomErrorResponse(c, envelope{"message": "Task already assigned"}, http.StatusConflict, fmt.Errorf("task already assigned"))
		return fmt.Errorf("task already assigned")
	}

	res, err := app.models.TaskRequests.CreateTaskRequest(taskWorkerId, taskId)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}
	RowsAffected, _ := res.RowsAffected()

	return c.JSON(http.StatusOK, envelope{"Rows Affected": RowsAffected})
}

func (app *application) removeTaskRequestHandler(c echo.Context) error {
	taskWorkerId, err1 := app.readIntParam(c, "userid")
	taskId, err2 := app.readIntParam(c, "taskid")
	if err1 != nil || err2 != nil {
		app.BadRequest(c, fmt.Errorf("error reading form data: %w", err1))
		return fmt.Errorf("error reading form data: %w", err1)
	}

	RequestedUser := app.contextGetUser(c)

	if int64(taskWorkerId) != RequestedUser.ID {
		app.CustomErrorResponse(c, envelope{"unauthorized": "You are not authorized to delete this task request"}, http.StatusUnauthorized, ErrUserUnauthorized)
		return ErrUserUnauthorized
	}

	res, err := app.models.TaskRequests.DeleteTaskRequest(taskWorkerId, taskId)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}
	RowsAffected, _ := res.RowsAffected()

	return c.JSON(http.StatusOK, envelope{"Rows Affected": RowsAffected})
}

func (app *application) GetQueryAuthorizeUser(c echo.Context) (int64, int64, int, error) {
	taskWorkerId, err1 := app.readIntParam(c, "userid")
	taskId, err2 := app.readIntParam(c, "taskid")
	if err1 != nil || err2 != nil {
		return 0, 0, 0, fmt.Errorf("error reading form data: %w", err1)
	}

	requestedUser := app.contextGetUser(c)

	taskOwner, taskVersion, err := app.models.Tasks.GetTaskForVerification(taskId)
	if err != nil {
		return 0, 0, 0, err
	}
	fmt.Println(taskOwner, requestedUser.ID, "taskOwner, taskVersion")

	if requestedUser.ID != taskOwner {
		return 0, 0, 0, ErrUserUnauthorized

	}

	return taskId, taskWorkerId, taskVersion, nil
}

// dormant
func (app *application) CheckoutTaskRequestHandler(c echo.Context) error {
	taskId, taskWorkerId, _, err := app.GetQueryAuthorizeUser(c)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	res, err := app.models.TaskRequests.GetCheckoutTaskRequest(taskWorkerId, taskId)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (app *application) approveTaskRequestHandler(c echo.Context) error {

	taskId, taskWorkerId, taskVersion, err := app.GetQueryAuthorizeUser(c)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	res, err := app.models.TaskRequests.ApproveTaskRequest(taskId, taskWorkerId, taskVersion)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	rowsAffected0, _ := res[0].RowsAffected()
	rowsAffected1, _ := res[1].LastInsertId()
	rowsAffected2, _ := res[2].RowsAffected()

	return c.JSON(http.StatusOK, envelope{
		"Task Rows Affected":            rowsAffected0,
		"Aproved Task Request Affected": rowsAffected1,
		"Task Request Rows Affected":    rowsAffected2,
	})
}

func (app *application) rejectTaskRequestHandler(c echo.Context) error {

	taskId, taskWorkerId, _, err := app.GetQueryAuthorizeUser(c)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	res, err := app.models.TaskRequests.RejectTaskRequest(taskId, taskWorkerId)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	rowsAffected, _ := res.RowsAffected()

	return c.JSON(http.StatusOK, envelope{"Rows Affected": rowsAffected})

}
