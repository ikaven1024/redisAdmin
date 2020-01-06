package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func errorHandle(c *gin.Context) {
	c.Next()

	if len(c.Errors) == 0 {
		return
	}

	err := c.Errors.Last().Err
	parsedError, ok := c.Errors.Last().Err.(*httpError)
	if !ok {
		parsedError = &httpError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	if parsedError.Code == 0 {
		parsedError.Code = http.StatusInternalServerError
	}

	c.JSON(parsedError.Code, result{
		Message: parsedError.Message,
		Detail: parsedError.Detail,
	})
}

type httpError struct {
	// Code will be used as http status code to response
	Code int
	// Message will be responded to client
	Message string
	// Detail will be printed in log
	Detail string
}

func NewBadRequestError(message string, detail string) error {
	return &httpError{Code: http.StatusBadRequest, Message: message, Detail: detail}
}

func NewServerError(message string, detail string) error {
	return &httpError{Code: http.StatusInternalServerError, Message: message, Detail: detail}
}

func (e httpError) Error() string {
	return fmt.Sprintf("%v: %v", e.Message, e.Detail)
}
