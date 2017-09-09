package api

import "octlink/ovs/modules/snat"

// AddSnat to add image by API
func AddSnat(paras *Paras) *Response {
	return nil
}

// ShowSnat by api
func ShowSnat(paras *Paras) *Response {
	return nil
}

// DeleteSnat to delete image
func DeleteSnat(paras *Paras) *Response {
	return &Response{
		Error: 0,
	}
}

// ShowAllSnats to display all images by condition
func ShowAllSnats(paras *Paras) *Response {
	return &Response{
		Data: snat.GetAllSnats(),
	}
}
