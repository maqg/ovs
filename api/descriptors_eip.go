package api

var eipDescriptors = Module{
	Name: "eip",
	Protos: map[string]Proto{

		"APIShowEips": {
			Name:    "查看所有EIP配置",
			handler: ShowEips,
			Paras:   []ProtoPara{},
		},

		"APIEipCreateEip": {
			Name:    "建立EIP配置",
			handler: CreateEip,
			Paras: []ProtoPara{
				{
					Name:    "privateMac",
					Type:    ParamTypeString,
					Desc:    "Mac Address of private nic",
					Default: ParamNotNull,
				},
				{
					Name:    "publicMac",
					Type:    ParamTypeString,
					Desc:    "Mac Address of public nic",
					Default: ParamNotNull,
				},
				{
					Name:    "vip",
					Type:    ParamTypeString,
					Desc:    "Virtual IP Address",
					Default: ParamNotNull,
				},
				{
					Name:    "guestIp",
					Type:    ParamTypeString,
					Desc:    "Guest IP Address",
					Default: ParamNotNull,
				},
			},
		},

		"APISyncEips": {
			Name:    "同步所有EIP配置",
			handler: SyncEips,
			Paras: []ProtoPara{
				{
					Name:    "eips",
					Type:    ParamTypeString,
					Desc:    "EIP Config in list []",
					Default: ParamNotNull,
				},
			},
		},

		"APIRemoveEips": {
			Name:    "删除所有EIP配置",
			handler: RemoveEips,
			Paras: []ProtoPara{
				{
					Name:    "eips",
					Type:    ParamTypeString,
					Desc:    "EIP Config in list []",
					Default: ParamNotNull,
				},
			},
		},

		"APIRemoveEip": {
			Name:    "删除EIP配置",
			handler: RemoveEip,
			Paras: []ProtoPara{
				{
					Name:    "privateMac",
					Type:    ParamTypeString,
					Desc:    "Mac Address of private nic",
					Default: ParamNotNull,
				},
				{
					Name:    "publicMac",
					Type:    ParamTypeString,
					Desc:    "Mac Address of public nic",
					Default: ParamNotNull,
				},
				{
					Name:    "vip",
					Type:    ParamTypeString,
					Desc:    "Virtual IP Address",
					Default: ParamNotNull,
				},
				{
					Name:    "guestIp",
					Type:    ParamTypeString,
					Desc:    "Guest IP Address",
					Default: ParamNotNull,
				},
			},
		},
	},
}
