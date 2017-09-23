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

		"APIShowAllDns": {
			Name:    "查看所有DNS",
			handler: ShowAllDns,
			Paras: []ProtoPara{
				{
					Name:    "start",
					Type:    ParamTypeInt,
					Desc:    "开始位置",
					Default: 0,
				},
				{
					Name:    "limit",
					Type:    "int",
					Desc:    "获取条目",
					Default: 15,
				},
			},
		},

		"APIRemoveDns": {
			Name:    "删除DNS",
			handler: DeleteDns,
			Paras: []ProtoPara{
				{
					Name:    "dnsAddress",
					Type:    ParamTypeString,
					Desc:    "dns server address",
					Default: ParamNotNull,
				},
			},
		},
	},
}
