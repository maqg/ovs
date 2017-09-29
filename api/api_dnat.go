package api

import "octlink/ovs/plugins"

// AddDnat to add dnat by API
func AddDnat(paras *Paras) *Response {

	dnat := &plugins.Dnat{
		VipPortStart:     paras.GetInt("vipPortStart"),
		VipPortEnd:       paras.GetInt("vipPortEnd"),
		PrivatePortStart: paras.GetInt("privatePortStart"),
		PrivatePortEnd:   paras.GetInt("privatePortEnd"),
		ProtocolType:     paras.Get("protocolType"),
		VipIp:            paras.Get("vipIp"),
		PrivateIp:        paras.Get("privateIp"),
		PrivateNicMac:    paras.Get("privateNicMac"),
		AllowedCidr:      paras.Get("allowedCidr"),
	}

	return &Response{
		Error: dnat.AddDnat(),
	}
}
