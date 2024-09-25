package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// bodyPtr is a POINTER to the body struct you are validating
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
