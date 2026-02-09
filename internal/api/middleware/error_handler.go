package middleware

import (
	"encoding/json"
	"errors"
	"net/http"

	myerrors "github.com/Coiiap5e/photographer/internal/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		var validationErr validator.ValidationErrors
		var unmarshalErr *json.UnmarshalTypeError
		var appErr *myerrors.AppError

		if errors.As(err, &validationErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "VALIDATION_ERROR",
				"message": "Invalid data provided: " + err.Error(),
			})
			return
		} else if errors.As(err, &unmarshalErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "INVALID_REQUEST",
				"message": "Invalid request format: " + err.Error(),
			})
			return
		} else if errors.As(err, &appErr) {
			switch appErr.Code {
			case myerrors.ErrCodeClientCreate:
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error":   "CREATE_CLIENT_ERROR",
					"message": "Failed to create client",
				})
			case myerrors.ErrCodeClientNotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"error":   "CLIENT_NOT_FOUND",
					"message": err.Error(),
				})
			case myerrors.ErrCodeShootCreate:
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error":   "CREATE_SHOOT_ERROR",
					"message": "Failed to create shoot",
				})
			case myerrors.ErrCodeShootUpdate:
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error":   "UPDATE_SHOOT_ERROR",
					"message": "Failed to update shoot",
				})
			case myerrors.ErrCodeShootNotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"error":   "SHOT_NOT_FOUND",
					"message": "Shoot not found",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "INTERNAL_ERROR",
					"message": "Internal server error",
				})
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "INTERNAL_ERROR",
				"message": "An unexpected server error occurred",
			})
		}
	}
}
