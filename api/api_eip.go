package api

import (
	"encoding/json"
	"octlink/ovs/plugins"
	"octlink/ovs/utils/merrors"
)

// CreateEip to add image by API
func CreateEip(paras *Paras) *Response {
	eip := &plugins.EipInfo{
		PrivateMac: paras.Get("privateMac"),
		PublicMac:  paras.Get("publicMac"),
		VipIP:      paras.Get("vip"),
		GuestIP:    paras.Get("guestIp"),
	}
	return &Response{
		Error: eip.CreateEip(),
	}
}

// RemoveEip by API
func RemoveEip(paras *Paras) *Response {
	eip := &plugins.EipInfo{
		PrivateMac: paras.Get("privateMac"),
		PublicMac:  paras.Get("publicMac"),
		VipIP:      paras.Get("vip"),
		GuestIP:    paras.Get("guestIp"),
	}
	return &Response{
		Error: eip.RemoveEip(),
	}
}

// RemoveEips by API
func RemoveEips(paras *Paras) *Response {
	return &Response{}
}

// SyncEips by API
func SyncEips(paras *Paras) *Response {

	eipsJSON := []byte(paras.Get("eips"))
	var eips []plugins.EipInfo

	err := json.Unmarshal(eipsJSON, &eips)
	if err != nil {
		return &Response{
			Error: merrors.ErrBadParas,
		}
	}

	eipsNew := make([]*plugins.EipInfo, len(eips))
	for i := range eips {
		eipsNew[i] = &eips[i]
	}

	logger.Debugf("eips paras:", eipsNew)

	return &Response{
		Error: plugins.SyncEips(eipsNew),
	}
}

// ShowEips by api
func ShowEips(paras *Paras) *Response {
	return &Response{
		Data: plugins.GetAllEips(),
	}
}

// ShowEip by api
func ShowEip(paras *Paras) *Response {

	eip, err := plugins.GetEip(paras.Get("privateMac"))

	return &Response{
		Data:  eip,
		Error: err,
	}

}
