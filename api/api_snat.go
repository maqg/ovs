package api

import "octlink/ovs/modules/snat"

// AddSnat to add image by API
func AddSnat(paras *Paras) *Response {

	sn := &snat.Snat{
		PrivateNicMac: paras.Get("privateNicMac"),
		PrivateNicIP:  paras.Get("privateNicIp"),
		SnatNetmask:   paras.Get("netmask"),
		PublicNicMac:  paras.Get("publicNicMac"),
		PublicIP:      paras.Get("publicIp"),
	}

	return &Response{
		Error: sn.Add(),
	}
}

// SyncSnat to add image by API
func SyncSnat(paras *Paras) *Response {

	sn := &snat.Snat{
		PrivateNicMac: paras.Get("privateNicMac"),
		PrivateNicIP:  paras.Get("privateNicIp"),
		SnatNetmask:   paras.Get("netmask"),
		PublicNicMac:  paras.Get("publicNicMac"),
		PublicIP:      paras.Get("publicIp"),
	}

	return &Response{
		Error: sn.Sync(),
	}
}

// ShowSnat by api
func ShowSnat(paras *Paras) *Response {

	privateIP := paras.Get("privateNicIp")
	netmask := paras.Get("netmask")

	return &Response{
		Data: snat.GetSnat(privateIP, netmask),
	}
}

// DeleteSnat to delete image
func DeleteSnat(paras *Paras) *Response {

	sn := snat.Snat{
		PrivateNicIP: paras.Get("privateNicIp"),
		SnatNetmask:  paras.Get("netmask"),
	}

	return &Response{
		Error: sn.Remove(),
	}
}

// ShowAllSnats to display all images by condition
func ShowAllSnats(paras *Paras) *Response {
	return &Response{
		Data: snat.GetAllSnats(),
	}
}
