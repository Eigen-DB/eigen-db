package hnsw_params

import (
	"eigen_db/api/utils"
	"eigen_db/cfg"
	"net/http"

	"github.com/gin-gonic/gin"
)

//type updateDimBody struct {
//	UpdatedDimensions int `json:"updatedDimensions" binding:"required"`
//}

func UpdateDimensions(config cfg.IConfig) func(*gin.Context) {
	return func(c *gin.Context) {
		utils.SendResponse(
			c,
			http.StatusOK,
			"Work in progress.",
			nil,
			nil,
		)

		//c.String(http.StatusOK, "Work in progesss.") // changing the vector space's dimensionality would invalidate all vectors stored within it. ill figure this one out later.
		//return

		/*var body updateDimBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		config.SetHNSWParamsDimensions(body.UpdatedDimensions)
		c.String(http.StatusOK, "Vector space dimension-count updated. Please restart the database for it to take effect.")*/
	}
}
