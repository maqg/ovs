package api

// dnatDescriptors for DNAT management by API
var dnatDescriptors = Module{
	Name: "dnat",
	Protos: map[string]Proto{

		"APIShowDnats": {
			Name:    "查看所有DNAT配置",
			handler: ShowDnats,
			Paras:   []ProtoPara{},
		},

		"APIShowDnat": {
			Name:    "查看DNAT配置",
			handler: ShowDnat,
			Paras: []ProtoPara{
				{
					Name:    "privateNicMac",
					Type:    ParamTypeString,
					Desc:    "Mac Address of private nic",
					Default: ParamNotNull,
				},
			},
		},

		"APIAddDnat": {
			Name:    "添加DNAT",
			handler: AddDnat,
			Paras: []ProtoPara{
				{
					Name:    "vipPortStart",
					Type:    ParamTypeInt,
					Desc:    "vip start port",
					Default: 22,
				},
				{
					Name:    "vipPortEnd",
					Type:    ParamTypeInt,
					Desc:    "vip end port",
					Default: 22,
				},
				{
					Name:    "privatePortStart",
					Type:    ParamTypeInt,
					Desc:    "private ip start port",
					Default: 22,
				},
				{
					Name:    "privatePortEnd",
					Type:    ParamTypeInt,
					Desc:    "private ip end port",
					Default: 22,
				},
				{
					Name:    "protocolType",
					Type:    ParamTypeString,
					Desc:    "prototol type",
					Default: ParamNotNull,
				},
				{
					Name:    "vipIp",
					Type:    ParamTypeString,
					Desc:    "Vip Ip Address",
					Default: ParamNotNull,
				},
				{
					Name:    "privateIp",
					Type:    ParamTypeString,
					Desc:    "Private IP Address",
					Default: ParamNotNull,
				},
				{
					Name:    "privateNicMac",
					Type:    ParamTypeString,
					Desc:    "Private Nic Mac Address",
					Default: ParamNotNull,
				},
				{
					Name:    "allowedCidr",
					Type:    ParamTypeString,
					Desc:    "allowed CIDR",
					Default: "",
				},
			},
		},
		"APISyncDnats": {
			Name:    "同步所有DNAT配置",
			handler: SyncDnats,
			Paras: []ProtoPara{
				{
					Name:    "dnats",
					Type:    ParamTypeString,
					Desc:    "DNAT Config in list []",
					Default: ParamNotNull,
				},
			},
		},
		"APIRemoveDnat": {
			Name:    "删除DNAT配置",
			handler: RemoveDnat,
			Paras: []ProtoPara{
				{
					Name:    "vipPortStart",
					Type:    ParamTypeInt,
					Desc:    "vip start port",
					Default: 22,
				},
				{
					Name:    "vipPortEnd",
					Type:    ParamTypeInt,
					Desc:    "vip end port",
					Default: 22,
				},
				{
					Name:    "privatePortStart",
					Type:    ParamTypeInt,
					Desc:    "private ip start port",
					Default: 22,
				},
				{
					Name:    "privatePortEnd",
					Type:    ParamTypeInt,
					Desc:    "private ip end port",
					Default: 22,
				},
				{
					Name:    "protocolType",
					Type:    ParamTypeString,
					Desc:    "prototol type",
					Default: ParamNotNull,
				},
				{
					Name:    "vipIp",
					Type:    ParamTypeString,
					Desc:    "Vip Ip Address",
					Default: ParamNotNull,
				},
				{
					Name:    "privateNicMac",
					Type:    ParamTypeString,
					Desc:    "Private Nic Mac Address",
					Default: ParamNotNull,
				},
			},
		},
		"APIRemoveDnats": {
			Name:    "删除所有DNAT配置",
			handler: RemoveDnats,
			Paras: []ProtoPara{
				{
					Name:    "dnats",
					Type:    ParamTypeString,
					Desc:    "DNAT Config in list []",
					Default: ParamNotNull,
				},
			},
		},
	},
}
