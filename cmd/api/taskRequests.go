package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *application) addNewTaskRequestHandler(c echo.Context) error {
	qs := c.Request().URL.Query()
	userId := int64(app.readIntQuery(qs, "userid", 0))
	taskId := int64(app.readIntQuery(qs, "taskid", 0))

	// ! commented for ease of testing, uncomment before deployment
	// requestedUser := app.contextGetUser(c)

	// if requestedUser.ID != userId {
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

	res, err := app.models.TaskRequests.CreateTaskRequest(userId, taskId)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}
	RowsAffected, _ := res.RowsAffected()

	return c.JSON(http.StatusOK, envelope{"Rows Affected": RowsAffected})
}

func (app *application) removeTaskRequestHandler(c echo.Context) error {
	qs := c.Request().URL.Query()
	userId := int64(app.readIntQuery(qs, "userid", 0))
	taskId := int64(app.readIntQuery(qs, "taskid", 0))

	RequestedUser := app.contextGetUser(c)

	if int64(userId) != RequestedUser.ID {
		app.CustomErrorResponse(c, envelope{"unauthorized": "You are not authorized to delete this task request"}, http.StatusUnauthorized, ErrUserUnauthorized)
		return ErrUserUnauthorized
	}

	res, err := app.models.TaskRequests.DeleteTaskRequest(userId, taskId)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}
	RowsAffected, _ := res.RowsAffected()

	return c.JSON(http.StatusOK, envelope{"Rows Affected": RowsAffected})
}

func getQueryAuthorizeUser(c echo.Context, app *application) (int64, int64, int, error) {
	qs := c.Request().URL.Query()
	userId := int64(app.readIntQuery(qs, "userid", 0))
	taskId := int64(app.readIntQuery(qs, "taskid", 0))

	requestedUser := app.contextGetUser(c)

	taskOwner, taskVersion, err := app.models.Tasks.GetTaskOwner(taskId)
	if err != nil {
		app.InternalServerError(c, err)
		return 0, 0, 0, err
	}

	if requestedUser.ID != taskOwner {
		app.CustomErrorResponse(c, envelope{"unauthorized": "You are not authorized to approve this task request"}, http.StatusUnauthorized, ErrUserUnauthorized)
		return 0, 0, 0, ErrUserUnauthorized

	}

	return taskId, userId, taskVersion, nil
}

func (app *application) approveTaskRequestHandler(c echo.Context) error {

	taskId, userId, taskVersion, err := getQueryAuthorizeUser(c, app)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	res, err := app.models.TaskRequests.ApproveTaskRequest(taskId, userId, taskVersion)
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

	taskId, userId, _, err := getQueryAuthorizeUser(c, app)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	res, err := app.models.TaskRequests.RejectTaskRequest(taskId, userId)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	rowsAffected, _ := res.RowsAffected()

	return c.JSON(http.StatusOK, envelope{"Rows Affected": rowsAffected})

}
