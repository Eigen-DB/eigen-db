package root_endpoints

import (
	"github.com/gin-gonic/gin"
)

type setConfigRequestBody struct {
	FieldName    string `json:"fieldName"`
	UpdatedValue any    `json:"updatedValue"`
}

func SetConfig(c *gin.Context) {
	//bodyBytes, err := io.ReadAll(c.Request.Body)
	//if err != nil {
	//	c.Error(err)
	//}
	//
	//body := &setConfigRequestBody{}
	//err = json.Unmarshal(bodyBytes, body)
	//if err != nil {
	//	c.Error(err)
	//}

	c.String(200, "Work in progress... The lead dev is sleeping (zzz)...")
}
