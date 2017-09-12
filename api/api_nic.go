package api

import (
	"octlink/ovs/modules/nic"
)

// ShowInterfaces by api
func ShowInterfaces(paras *Paras) *Response {
	return &Response{
		Data: nic.GetNics(),
	}
}
