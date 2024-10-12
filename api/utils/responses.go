package utils

import (
	"github.com/gin-gonic/gin"
)

// A response error.
//
// Response errors have the following schema:
//
//	{
//		"code": "SOME_CODE",
//		"description": "some description..."
//	}
type resErr struct {
	code        string
	description any
}

// Sends a structured HTTP response, given specified values.
//
// The schema for responses in EigenDB's API is the following:
//
// ("data" and "error" are optional fields)
//
//	{
//		"status": 200,
//		"message": "Request successful",
//		"data": {
//			"id": 1,
//			"name": "John Doe"
//		},
//		"error": {
//			"code": "RESOURCE_NOT_FOUND",
//			"description": "The requested resource could not be located"
//		}
//	}
func SendResponse(ginCtx *gin.Context, statusCode int, message string, data map[string]any, err *resErr) {
	res := gin.H{
		"status":  statusCode,
		"message": message,
	}

	if data != nil {
		res["data"] = data
	}
	if err != nil {
		res["error"] = map[string]any{
			"code":        err.code,
			"description": err.description,
		}
	}

	ginCtx.JSON(statusCode, res) // sends the HTTP response
}

// Creates an instance of ResErr given specified values
//
// Returns a pointer to the created instance.
func CreateError(code string, desc any) *resErr {
	return &resErr{
		code:        code,
		description: desc,
	}
}
