package controllers

import (
	"bookingservice/utils"

	"github.com/gobuffalo/buffalo"
)

func HealthHandler(c buffalo.Context) error {
	return utils.DoResponse(c, 200, 0, "Healthy", "")
}
