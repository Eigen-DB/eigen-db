package indexes

import (
	"eigen_db/api/utils"
	"eigen_db/index_mgr"

	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	idx_list, err := index_mgr.GetIndexMgr().ListIndexes()
	if err != nil {
		utils.SendResponse(
			c,
			500,
			"An error occured while listing indexes.",
			nil,
			utils.CreateError("INDEXES_NOT_LISTED", err.Error()),
		)
		return
	}

	utils.SendResponse(
		c,
		200,
		"Indexes listed successfully.",
		map[string]any{
			"indexes": idx_list,
		},
		nil,
	)
}
