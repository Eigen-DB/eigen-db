package utils

import (
	"github.com/gin-gonic/gin"
)

type ResErr struct {
	code        string
	description any
}

/*
Response schema:

"data" and "error" are optional fields

	{
		"status": 200,
		"message": "Request successful",
		"data": {
			"id": 1,
			"name": "John Doe"
		},
		"error": {
			"code": "RESOURCE_NOT_FOUND",
			"description": "The requested resource could not be located"
		}
	}
*/
func SendResponse(ginCtx *gin.Context, statusCode int, message string, data map[string]any, err *ResErr) {
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

	ginCtx.JSON(statusCode, res)
}

func CreateError(code string, desc any) *ResErr {
	return &ResErr{
		code:        code,
		description: desc,
	}
}
