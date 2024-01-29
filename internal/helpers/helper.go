package helpers

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

func ConvertToInt64(str string) (int64, error) {
	integer, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return integer, nil
}

func ParamToInt64(c echo.Context, str string) (int64, error) {
	param := c.Param("id")
	id, err := ConvertToInt64(param)
	if err != nil {
		return 0, err
	}
	return id, err
}
