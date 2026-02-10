package middleware

import (
	"log"
	"net/http"

	"bookingservice/exceptions"
	"bookingservice/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gobuffalo/buffalo"
)

func ErrorHandler() buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {

			defer func() {
				if err := recover(); err != nil {
					log.Printf("PANIC: %+v\n", err)
					switch e := err.(type) {
					case *exceptions.AuthorizationError:
						_ = utils.DoResponse(c, http.StatusForbidden, 1401, e.Error(), "")
						return
					case *exceptions.UserNotFoundError:
						_ = utils.DoResponse(c, http.StatusNotFound, 1404, e.Error(), "")
					case *exceptions.WrongPasswordError:
						_ = utils.DoResponse(c, http.StatusUnauthorized, 1401, e.Error(), "")
					case *exceptions.InvalideRequestError:
						_ = utils.DoResponse(c, http.StatusBadRequest, 1400, e.Error(), "")
					case validator.ValidationErrors:
						_ = utils.DoResponse(c, http.StatusBadRequest, 1400, e.Error(), "")
					default:
						_ = utils.DoResponse(c, http.StatusInternalServerError, 1500, "Internal Server Error", "")
					}
				}
			}()

			err := next(c)
			if err != nil {
				switch e := err.(type) {
				case *exceptions.InvalideRequestError:
					return utils.DoResponse(c, http.StatusBadRequest, 1400, e.Error(), "")
				default:
					return utils.DoResponse(c, http.StatusInternalServerError, 1500, "Internal Server Error", "")
				}
			}
			return nil
		}
	}
}
