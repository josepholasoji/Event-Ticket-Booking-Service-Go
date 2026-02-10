package utils

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"

	"bookingservice/dtos/responses"
)

var r = render.New(render.Options{})

func DoResponse(c buffalo.Context, httpStatusCode int, responseCode int, message string, data string) error {
	c.Response().Header().Set("Content-Type", "application/json")
	return c.Render(httpStatusCode, r.JSON(responses.Response{
		Code:    responseCode,
		Message: message,
		Data:    data,
	}))
}
