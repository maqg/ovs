package api

import "octlink/ovs/plugins"

// ShowInterfaces by api
func ShowInterfaces(paras *Paras) *Response {
	return &Response{
		Data: plugins.GetNics(),
	}
}
