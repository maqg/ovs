package api

import "octlink/ovs/plugins"

// ShowInterfaces by api
func ShowInterfaces(paras *Paras) *Response {
	return &Response{
		Data: plugins.GetNics(),
	}
}

// SetInterface by api
func SetInterface(paras *Paras) *Response {

	ifInfo := &plugins.IfInfo{
		IP:      paras.Get("ip"),
		Mac:     paras.Get("mac"),
		Netmask: paras.Get("netmask"),
	}

	return &Response{
		Data: ifInfo.ConfigureNic(),
	}
}

// RemoveInterface by api
func RemoveInterface(paras *Paras) *Response {
	ifInfo := &plugins.IfInfo{
		Mac: paras.Get("mac"),
	}
	return &Response{
		Data: ifInfo.RemoveNic(),
	}
}
