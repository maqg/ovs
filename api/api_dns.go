package api

import "octlink/ovs/plugins"

// AddDns to add dns
func AddDns(paras *Paras) *Response {

	dns := &plugins.Dns{
		DnsAddress:   paras.Get("dnsAddress"),
		PublicNicMac: paras.Get("publicNicMac"),
	}

	return &Response{
		Error: dns.AddDns(),
	}
}

// DeleteDns for delete dns
func DeleteDns(paras *Paras) *Response {

	dns := &plugins.Dns{
		DnsAddress: paras.Get("dnsAddress"),
	}

	return &Response{
		Error: dns.DeleteDns(),
	}
}
