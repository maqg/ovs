package api

import "octlink/ovs/modules/systemconfig"

// ShowSystemConfig to add image by API
func ShowSystemConfig(paras *Paras) *Response {
	return &Response{
		Data: systemconfig.GetSystemConfig(),
	}
}
