package api

import "octlink/ovs/modules/snat"
import "octlink/ovs/utils/merrors"

// AddSnat to add image by API
func AddSnat(paras *Paras) *Response {
	sn := snat.GetSnat(paras.Get("privateNicMac"))
	if sn != nil {
		return &Response{
			Error: merrors.ErrSegmentAlreadyExist,
		}
	}

	sn = &snat.Snat{
		PrivateNicMac: paras.Get("privateNicMac"),
	}

	return &Response{
		Error: sn.Add(),
	}
}

// ShowSnat by api
func ShowSnat(paras *Paras) *Response {
	return &Response{}
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
