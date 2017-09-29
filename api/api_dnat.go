package api

import (
	"encoding/json"
	"octlink/ovs/plugins"
	"octlink/ovs/utils/merrors"
)

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

// RemoveDnat to remove dnat
func RemoveDnat(paras *Paras) *Response {

	dnat := plugins.Dnat{
		VipPortStart:     paras.GetInt("vipPortStart"),
		VipPortEnd:       paras.GetInt("vipPortEnd"),
		PrivatePortStart: paras.GetInt("privatePortStart"),
		PrivatePortEnd:   paras.GetInt("privatePortEnd"),
		ProtocolType:     paras.Get("protocolType"),
		VipIp:            paras.Get("vipIp"),
		PrivateNicMac:    paras.Get("privateNicMac"),
	}

	return &Response{
		Error: dnat.RemoveDnat(),
	}
}

// RemoveDnats by API
func RemoveDnats(paras *Paras) *Response {
	dnatsJSON := []byte(paras.Get("dnats"))
	var dnats []plugins.Dnat

	err := json.Unmarshal(dnatsJSON, &dnats)
	if err != nil {
		return &Response{
			Error: merrors.ErrBadParas,
		}
	}

	dnatsNew := make([]*plugins.Dnat, len(dnats))
	for i := range dnats {
		dnatsNew[i] = &dnats[i]
	}

	return &Response{
		Error: plugins.RemoveDnats(dnatsNew),
	}
}

// SyncDnats by API
func SyncDnats(paras *Paras) *Response {

	dnatsJSON := []byte(paras.Get("dnats"))
	var dnats []plugins.Dnat

	err := json.Unmarshal(dnatsJSON, &dnats)
	if err != nil {
		return &Response{
			Error: merrors.ErrBadParas,
		}
	}

	dnatsNew := make([]*plugins.Dnat, len(dnats))
	for i := range dnats {
		dnatsNew[i] = &dnats[i]
	}

	return &Response{
		Error: plugins.SyncDnats(dnatsNew),
	}
}

// ShowDnats by api
func ShowDnats(paras *Paras) *Response {
	return &Response{
		Data: plugins.GetAllDnats(),
	}
}

// ShowDnat by api
func ShowDnat(paras *Paras) *Response {

	dnat, err := plugins.GetDnat(paras.Get("privateNicMac"))

	return &Response{
		Data:  dnat,
		Error: err,
	}

}
