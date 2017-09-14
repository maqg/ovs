package api

import (
	"octlink/ovs/plugins"
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
		Data: eip,
	}
}

// RemoveEip by API
func RemoveEip(paras *Paras) *Response {
	return &Response{}
}

// RemoveEips by API
func RemoveEips(paras *Paras) *Response {
	return &Response{}
}

// SyncEip by API
func SyncEip(paras *Paras) *Response {
	return &Response{}
}

// SyncEips by API
func SyncEips(paras *Paras) *Response {
	return &Response{}
}

// ShowEips by api
func ShowEips(paras *Paras) *Response {
	return &Response{}
}

// ShowEip by api
func ShowEip(paras *Paras) *Response {
	return &Response{}
}
