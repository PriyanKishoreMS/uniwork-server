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
	e.Use(middleware.CORS())
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(app.rateLimit())

	e.GET("/health", app.HealthCheckHandler)
	e.POST("/user", app.registerUserHandler)
	e.POST("/login", app.loginUserHandler)
	e.POST("/security/refreshtoken", app.refreshTokenHandler)
	e.Static("/public", uploadDir)
	e.GET("/pay", app.checkPayementHander)

	college := e.Group("/college", app.authenticate())
	{
		college.GET("", app.listAllCollegesHandler)
		college.GET("/:id", app.getCollegeHandler)
		college.POST("", app.createCollegeHandler)
		college.PATCH("/:id", app.updateCollegeHandler)
		college.DELETE("/:id", app.deleteCollegeHandler)
	}

	user := e.Group("/user", app.authenticate())
	{
		user.GET("/college/:id", app.listAllUsersInCollegeHandler)
		user.GET("/:id", app.getUserHandler)
		user.GET("", app.getRequestedUserHandler)
		user.PATCH("", app.updateUserHandler)
		user.DELETE("", app.deleteUserHandler)
	}

	service := e.Group("/task", app.authenticate())
	{
		service.GET("", app.listAllTasksHandler)
		service.GET("/user/:uid", app.listAllTasksOfUserHandler)
		service.GET("/worker/:uid", app.listAllTasksOfUserHandler)
		service.GET("/:id", app.getTaskHandler)
		service.POST("", app.addNewTaskHandler)
		service.DELETE("/:id", app.deleteTaskHandler)

		service.POST("/request/:taskid/:userid", app.addNewTaskRequestHandler)
		service.PATCH("/request/approve/:taskid/:userid", app.approveTaskRequestHandler)
		service.PATCH("/request/reject/:taskid/:userid", app.rejectTaskRequestHandler)
		service.DELETE("/request/:taskid/:userid", app.removeTaskRequestHandler)
	}

	return e
}
