package api

import (
	"octlink/ovs/utils"
	"octlink/ovs/utils/httpresponse"
	"octlink/ovs/utils/merrors"
	"reflect"

	"github.com/gin-gonic/gin"
)

// Response structure
type Response struct {
	Error    int         `json:"error"`
	ErrorLog string      `json:"errorLog"`
	Data     interface{} `json:"data"`
	Total    int         `json:"total"`
	Count    int         `json:"count"`
}

/*
{
	"module": "octlink.virtualrouter.center.host.APIAddHost",
	"paras": {
		"ip": "kk",
		"account": "root",
		"password": ""
	},
	"async": false,
}
*/
type inputParas struct {
	Module string
	API    string
	Paras  map[string]interface{}
	Async  bool
}

// Paras of API
type Paras struct {
	Proto   *Proto
	InParas *inputParas
}

// Get paras from Paras structure
func (p *Paras) Get(name string) string {
	if v := p.InParas.Paras[name]; v != nil {
		return v.(string)
	}
	return ""
}

// GetBoolean get boolean para, name can be false,true,0,1
func (p *Paras) GetBoolean(name string) bool {
	raw := p.InParas.Paras[name]
	if reflect.TypeOf(raw).Kind() != reflect.Bool {
		return utils.NumberToInt(raw) != 0
	}
	return p.InParas.Paras[name].(bool)
}

// GetInt paras from Paras structure
func (p *Paras) GetInt(name string) int {
	return utils.StringToInt(p.Get(name))
}

// GetInt64 paras from Paras structure
func (p *Paras) GetInt64(name string) int64 {
	return utils.StringToInt64(p.Get(name))
}

// Test for api test page
func (api *API) Test(c *gin.Context) {
	httpresponse.Ok(c, "Api Server is Running")
}

// GServices for api service management
var GServices map[string]*Service

// Service of API
type Service struct {
	Name    string                 `json:"name"`
	Handler func(*Paras) *Response `json:"handler"`
}

// GetService for api
func GetService(key string) *Service {
	service, ok := GServices[key]
	if !ok {
		logger.Errorf("no service for %s found\n", key)
		return nil
	}

	return service
}

func getParas(c *gin.Context) (*Paras, int) {

	var apiParas = new(Paras)

	c.BindJSON(&apiParas.InParas)

	logger.Debugf("got api %s\n", apiParas.InParas.API)

	if apiParas.InParas.API == "" {
		logger.Errorf("got null api\n")
		return nil, merrors.ErrNoSuchAPI
	}

	proto := FindProto(apiParas.InParas.API)
	if proto == nil {
		logger.Errorf("no api proto found for %s\n",
			apiParas.InParas.API)
		return nil, merrors.ErrNoSuchAPI
	}

	apiParas.Proto = proto

	return apiParas, 0
}

func checkParas(apiParas *Paras) (int, string) {

	protoParas := apiParas.Proto.Paras

	for i := 0; i < len(protoParas); i++ {

		protoParam := protoParas[i]
		inParam := apiParas.InParas.Paras[protoParam.Name]

		// if paras have default value and no input sepecified, set a default value
		if protoParam.Default != ParamNotNull && inParam == nil {
			apiParas.InParas.Paras[protoParam.Name] = protoParam.Default
		}

		logger.Debugf("param:%s, default:%s, value:%s\n", protoParam.Name,
			protoParam.Default, inParam)

		if protoParam.Default == ParamNotNull && inParam.(string) == "" {
			errorMsg := "paras \"" + protoParam.Name + "\" must be specified"
			return merrors.ErrNotEnoughParas, errorMsg
		}
	}

	return merrors.ErrSuccess, ""
}

// Dispatch api request
func (api *API) Dispatch(c *gin.Context) {

	logger.Debugf("got api request\n")

	paras, err := getParas(c)
	if paras == nil {
		logger.Errorf("No match proto found\n")
		httpresponse.Error(c, err, nil)
		return
	}

	service := GetService(paras.InParas.API)
	if service == nil {
		logger.Errorf("No match service found\n")
		httpresponse.Error(c, merrors.ErrNoSuchAPI, paras.InParas.API)
		return
	}

	ret, msg := checkParas(paras)
	if ret != merrors.ErrSuccess {
		logger.Errorf("Not Enough Paras\n")
		httpresponse.Error(c, merrors.ErrNotEnoughParas, msg)
		return
	}

	resp := service.Handler(paras)

	if resp.Error == 0 {
		httpresponse.Ok(c, resp.Data)
	} else {
		httpresponse.Error(c, resp.Error, resp.ErrorLog)
	}
}
