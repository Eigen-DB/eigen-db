package api

import (
	vio "eigen_db/vector_io"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
)

type searchRequestBody struct {
	QueryVectorId vio.VectorId `json:"queryVectorId"`
	K             uint32       `json:"k"`
}

func Search(c *gin.Context) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Error(err)
	}

	body := &searchRequestBody{}
	err = json.Unmarshal(bodyBytes, body)
	if err != nil {
		c.Error(err)
	}

	nnIds := vio.SimilaritySearch(body.QueryVectorId, body.K)

	jsonResponse, err := json.Marshal(nnIds)
	if err != nil {
		c.Error(err)
	}

	c.String(200, string(jsonResponse))
}
