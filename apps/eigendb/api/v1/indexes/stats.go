package indexes

import (
	"eigen_db/api/utils"
	"eigen_db/index_mgr"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Stats(c *gin.Context) {
	indexName := c.Param("index")
	idx, err := index_mgr.GetIndexMgr().GetIndex(indexName)
	if err != nil {
		utils.SendResponse(
			c,
			500,
			"An error occured while fetching the index.",
			nil,
			utils.CreateError("INDEX_NOT_FETCHED", err.Error()),
		)
		return
	}

	stats := map[string]any{
		"index_name": indexName,
		"dimensions": idx.Dimensions,
		"metric":     idx.Metric,
	}
	utils.SendResponse(
		c,
		200,
		fmt.Sprintf("Stats for index '%s' fetched successfully.", indexName),
		stats,
		nil,
	)
}
