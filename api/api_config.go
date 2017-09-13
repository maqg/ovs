package api

import "octlink/ovs/plugins"

// ShowSystemConfig to add image by API
func ShowSystemConfig(paras *Paras) *Response {
	return &Response{
		Data: plugins.GetSystemConfig(),
	}
}
