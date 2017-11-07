package api

// dnsDescriptors for VIP management by API
var vipDescriptors = Module{
	Name: "vip",
	Protos: map[string]Proto{

		"APIAddVip": {
			Name:    "添加VIP",
			handler: AddVip,
			Paras: []ProtoPara{
				{
					Name:    "ip",
					Type:    ParamTypeString,
					Desc:    "Virtual Ip",
					Default: ParamNotNull,
				},
				{
					Name:    "netmask",
					Type:    ParamTypeString,
					Desc:    "Virtual Ip Netmask",
					Default: ParamNotNull,
				},
				{
					Name:    "ownerEthernetMac",
					Type:    ParamTypeString,
					Desc:    "Vip Owner Ethernet Mac",
					Default: ParamNotNull,
				},
			},
		},
		"APIRemoveVip": {
			Name:    "删除VIP",
			handler: DeleteVip,
			Paras: []ProtoPara{
				{
					Name:    "ip",
					Type:    ParamTypeString,
					Desc:    "Virtual Ip",
					Default: ParamNotNull,
				},
				{
					Name:    "netmask",
					Type:    ParamTypeString,
					Desc:    "Virtual Ip Netmask",
					Default: ParamNotNull,
				},
				{
					Name:    "ownerEthernetMac",
					Type:    ParamTypeString,
					Desc:    "Vip Owner Ethernet Mac",
					Default: ParamNotNull,
				},
			},
		},

		"APISyncVips": {
			Name:    "同步所有VIP",
			handler: SyncVips,
			Paras: []ProtoPara{
				{
					Name:    "vips",
					Type:    ParamTypeString,
					Desc:    "VIP Config in list []",
					Default: ParamNotNull,
				},
			},
		},
	},
}
