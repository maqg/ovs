package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoadTestPage to load api test page
func (api *API) LoadTestPage(c *gin.Context) {
	apiModules, _ := json.Marshal(GAPIConfig.Modules)
	c.HTML(http.StatusOK, "apitest.html",
		gin.H{
			"TESTTITLE": "Mirage",
			"APICONFIG": string(apiModules),
		})
}
