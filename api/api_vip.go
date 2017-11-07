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

func SyncVips(paras *Paras) *Response {

	vipsJSON := []byte(paras.Get("vips"))
	var vips []plugins.Vip

	err := json.Unmarshal(vipsJSON, &vips)
	if err != nil {
		return &Response{
			Error: merrors.ErrBadParas,
		}
	}

	vipsNew := make([]*plugins.Vip, len(vips))
	for i := range vips {
		vipsNew[i] = &vips[i]
	}

	return &Response{
		Error: plugins.SyncVips(vipsNew),
	}
}
