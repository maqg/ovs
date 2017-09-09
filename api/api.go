package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"octlink/ovs/utils/octlog"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// ParamTypeString string type param
	ParamTypeString = "string"

	// ParamTypeInt int type param
	ParamTypeInt = "int"

	// ParamTypeListInt param type of list int
	ParamTypeListInt = "listint"

	// ParamTypeListString param type of list string
	ParamTypeListString = "liststring"

	// ParamTypeBoolean boolean type param
	ParamTypeBoolean = "boolean"

	// ParamNotNull not null param
	ParamNotNull = "NotNull"

	// APIPrefixCenter API Prefix of Center API
	APIPrefixCenter = "octlink.virtualrouter.v5"
)

var logger *octlog.LogConfig

// GAPIConfig for api config
var GAPIConfig Config

// API structure
type API struct {
	Name   string
	Config *Config
}

// Config structure
type Config struct {
	Modules map[string]Module `json:"modules"`
}

// ProtoPara of API
type ProtoPara struct {
	Name    string      `json:"name"`
	Default interface{} `json:"default"`
	Type    string      `json:"type"`
	Desc    string      `json:"desc"`
}

// Proto API proto structure
type Proto struct {
	Name    string      `json:"name"`
	Key     string      `json:"key"`
	Paras   []ProtoPara `json:"paras"`
	handler func(*Paras) *Response
}

// InitLog to init api log config
func InitLog(level int) {
	logger = octlog.InitLogConfig("api.log", level)
}

// Module of API
type Module struct {
	Name   string           `json:"name"`
	Protos map[string]Proto `json:"protos"`
}

// FindProto for by api key like xx.xxx.xxx.xx
func FindProto(api string) *Proto {

	segments := strings.Split(api, ".")
	moduleName := segments[3]
	apiKey := segments[4]

	if moduleName == "" || apiKey == "" {
		fmt.Printf("got bad api key %s\n", api)
		return nil
	}

	module, ok := GAPIConfig.Modules[moduleName]
	if !ok {
		fmt.Printf("no module exist for %s\n", moduleName)
		return nil
	}

	proto, ok := module.Protos[apiKey]
	if !ok {
		fmt.Printf("no proto exist for %s\n", apiKey)
		return nil
	}

	return &proto
}

// LoadTestPage to load api test page
func (api *API) LoadTestPage(c *gin.Context) {
	apiModules, _ := json.Marshal(GAPIConfig.Modules)
	c.HTML(http.StatusOK, "apitest.html",
		gin.H{
			"TESTTITLE": "Mirage",
			"APICONFIG": string(apiModules),
		})
}
