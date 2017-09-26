package api

import "octlink/ovs/plugins"

// AddVip to add image by API
func AddVip(paras *Paras) *Response {
	vip := &plugins.Vip{
		Ip:               paras.Get("ip"),
		Netmask:          paras.Get("netmask"),
		OwnerEthernetMac: paras.Get("ownerEthernetMac"),
	}

	return &Response{
		Error: vip.AddVip(),
	}
}

// DeleteVip to delete vip
func DeleteVip(paras *Paras) *Response {
	vip := &plugins.Vip{
		Ip:               paras.Get("ip"),
		Netmask:          paras.Get("netmask"),
		OwnerEthernetMac: paras.Get("ownerEthernetMac"),
	}

	return &Response{
		Error: vip.DeleteVip(),
	}
}
