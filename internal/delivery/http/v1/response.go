package v1

import "github.com/gin-gonic/gin"

//type dataResponse struct {
//	Data  interface{} `json:"data"`
//	Count int64       `json:"count"`
//}
//
//type idResponse struct {
//	ID interface{} `json:"id"`
//}

type response struct {
	Message string `json:"message"`
}

// ErrorResponseMessageOverride aborts the context and sends a properly formatted response with the specified message
// TODO: log error data
func ErrorResponseMessageOverride(c *gin.Context, statusCode int, _ error, message string) {
	//logger.Error(err)
	c.AbortWithStatusJSON(statusCode, response{message})
}

// ErrorResponse aborts the context and sends a properly formatted response with the message from supplied error
func ErrorResponse(c *gin.Context, statusCode int, err error) {
	//logger.Error(err)
	c.AbortWithStatusJSON(statusCode, response{err.Error()})
}

// ErrorResponseString aborts the context and sends a properly formatted response with the message from supplied string
func ErrorResponseString(c *gin.Context, statusCode int, message string) {
	//logger.Error(message)
	c.AbortWithStatusJSON(statusCode, response{message})
}
