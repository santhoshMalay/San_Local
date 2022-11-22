package utils

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
)

//type ValidationErrors = validation.Errors

//type dataResponse struct {
//	Data  interface{} `json:"data"`
//	Count int64       `json:"count"`
//}
//
//type idResponse struct {
//	ID interface{} `json:"id"`
//}

// Response corresponds to RFC 7807 Problem Details object
// https://www.rfc-editor.org/rfc/rfc7807
type Response struct {
	Title  string `json:"title"`
	Status int    `json:"status"`
}

type ValidationError struct {
	Title            string            `json:"title" example:"invalid request parameters"`
	Status           int               `json:"status" example:"400"`
	ValidationErrors validation.Errors `json:"validation_errors,omitempty"`
}

// ErrorResponseMessageOverride aborts the context and sends a properly formatted response with the specified message
// TODO: log error data
func ErrorResponseMessageOverride(c *gin.Context, statusCode int, _ error, message string) {
	//logger.Error(err)
	c.AbortWithStatusJSON(statusCode, Response{
		Title:  message,
		Status: statusCode,
	})
}

// ErrorResponse aborts the context and sends a properly formatted response with the message from supplied error
func ErrorResponse(c *gin.Context, statusCode int, err error) {
	//logger.Error(err)
	c.AbortWithStatusJSON(statusCode, Response{
		Title:  err.Error(),
		Status: statusCode,
	})
}

// ErrorResponseString aborts the context and sends a properly formatted response with the message from supplied string
func ErrorResponseString(c *gin.Context, statusCode int, message string) {
	//logger.Error(message)
	c.AbortWithStatusJSON(statusCode, Response{
		Title:  message,
		Status: statusCode,
	})
}

func ValidationErrorResponse(c *gin.Context, validationErrors validation.Errors) {
	//logger.Error(err)
	c.AbortWithStatusJSON(http.StatusBadRequest, ValidationError{
		Title:            "invalid request parameters",
		Status:           http.StatusBadRequest,
		ValidationErrors: validationErrors,
	})
}
