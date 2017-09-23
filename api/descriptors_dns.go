package api

// dnsDescriptors for DNS management by API
var dnsDescriptors = Module{
	Name: "dns",
	Protos: map[string]Proto{

		"APIAddDns": {
			Name:    "添加DNS",
			handler: AddDns,
			Paras: []ProtoPara{
				{
					Name:    "dnsAddress",
					Type:    ParamTypeString,
					Desc:    "dns server address",
					Default: ParamNotNull,
				},
				{
					Name:    "publicNicMac",
					Type:    ParamTypeString,
					Desc:    "Public Nic Mac Address",
					Default: ParamNotNull,
				},
			},
		},
	},
}
