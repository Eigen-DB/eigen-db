package indexes

import (
	"eigen_db/api/utils"
	"eigen_db/index_mgr"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Delete(c *gin.Context) {
	indexName := c.Param("index")
	if err := index_mgr.GetIndexMgr().DeleteIndex(indexName); err != nil {
		utils.SendResponse(
			c,
			500,
			"An error occured while deleting the index.",
			nil,
			utils.CreateError("INDEX_NOT_DELETED", err.Error()),
		)
		return
	}

	utils.SendResponse(
		c,
		200,
		fmt.Sprintf("Index '%s' successfully deleted.", indexName),
		nil,
		nil,
	)
}
