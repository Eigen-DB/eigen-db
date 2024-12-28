package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Validates a request body.
//
// The request body is stored in the ginCtx and bodyPtr
// is a POINTER to an instance of the body struct for the
// endpoint you are validating.
//
// If body is invalid, the following 404 response is returned and an error is returned:
//
//	{
//		"status": 400,
//		"message": "Bad request",
//		"error": {
//			"code": "INVALID_REQUEST_BODY",
//			"description": "The body you provided in your request is invalid."
//		}
//	}
//
// Example usage:
//
//	type addItemBody struct {
//		ItemName string `json:"itemName" binding:"required"`
//		ItemPrice float32 `json:"itemPrice" binding:"required"`
//	}
//
//	func AddItem(c *gin.Context) {
//		var body addItemBody
//		if err := utils.ValidateBody(c, &body); err != nil { // validates the request body
//			return
//		}
//	}
func ValidateBody(ginCtx *gin.Context, bodyPtr any) error {
	if err := ginCtx.ShouldBindJSON(bodyPtr); err != nil {
		SendResponse(
			ginCtx,
			http.StatusBadRequest,
			"Bad request",
			nil,
			CreateError("INVALID_REQUEST_BODY", "The body you provided in your request is invalid."),
		)
		return err
	}
	return nil
}
