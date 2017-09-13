package api

// snatDescriptors for SNAT management by API
var snatDescriptors = Module{
	Name: "snat",
	Protos: map[string]Proto{
		"APIAddSnat": {
			Name:    "添加SNAT",
			handler: AddSnat,
			Paras: []ProtoPara{
				{
					Name:    "privateNicMac",
					Type:    ParamTypeString,
					Desc:    "Private Nic Mac Address",
					Default: ParamNotNull,
				},
				{
					Name:    "privateNicIp",
					Type:    ParamTypeString,
					Desc:    "Private Nic IP Address",
					Default: "0.0.0.0",
				},
				{
					Name:    "netmask",
					Type:    ParamTypeString,
					Desc:    "Private Netmask Address",
					Default: "0.0.0.0",
				},
				{
					Name:    "publicNicMac",
					Type:    ParamTypeString,
					Desc:    "Public Nic Mac Address",
					Default: ParamNotNull,
				},
				{
					Name:    "publicIp",
					Type:    ParamTypeString,
					Desc:    "Public IP Address",
					Default: ParamNotNull,
				},
			},
		},

		"APIShowSnat": {
			Name:    "查看单个SNAT",
			handler: ShowSnat,
			Paras: []ProtoPara{
				{
					Name:    "privateNicMac",
					Type:    ParamTypeString,
					Desc:    "Private Nic Mac Address",
					Default: ParamNotNull,
				},
			},
		},

		"APIShowAllSnat": {
			Name:    "查看所有SNAT",
			handler: ShowAllSnats,
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

		"APIRemoveSnat": {
			Name:    "删除SNAT",
			handler: DeleteSnat,
			Paras: []ProtoPara{
				{
					Name:    "privateNicMac",
					Type:    ParamTypeString,
					Desc:    "Private Nic Mac Address",
					Default: ParamNotNull,
				},
			},
		},
	},
}
