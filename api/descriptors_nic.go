package api

var nicDescriptors = Module{
	Name: "nic",
	Protos: map[string]Proto{
		"APIShowInterfaces": {
			Name:    "查看接口信息",
			handler: ShowInterfaces,
			Paras:   []ProtoPara{},
		},
		"APISetInterface": {
			Name:    "设置接口信息",
			handler: SetInterface,
			Paras: []ProtoPara{
				{
					Name:    "mac",
					Type:    ParamTypeString,
					Desc:    "Mac Address of this nic",
					Default: ParamNotNull,
				},
				{
					Name:    "ip",
					Type:    ParamTypeString,
					Desc:    "IP Address of this nic",
					Default: ParamNotNull,
				},
				{
					Name:    "netmask",
					Type:    ParamTypeString,
					Desc:    "Netmask of Address",
					Default: ParamNotNull,
				},
			},
		},
		"APIRemoveInterface": {
			Name:    "删除接口配置",
			handler: RemoveInterface,
			Paras: []ProtoPara{
				{
					Name:    "mac",
					Type:    ParamTypeString,
					Desc:    "Mac Address of this nic",
					Default: ParamNotNull,
				},
			},
		},
	},
}
