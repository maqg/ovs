package api

// SnatDescriptors for SNAT management by API
var SnatDescriptors = Module{
	Name: "snat",
	Protos: map[string]Proto{
		"APIAddSnat": {
			Name:    "添加SNAT",
			handler: AddSnat,
			Paras: []ProtoPara{
				{
					Name:    "name",
					Type:    ParamTypeString,
					Desc:    "Image Name",
					Default: ParamNotNull,
				},
			},
		},

		"APIShowSnat": {
			Name:    "查看单个SNAT",
			handler: ShowSnat,
			Paras: []ProtoPara{
				{
					Name:    "id",
					Type:    ParamTypeString,
					Desc:    "Image Id",
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
					Name:    "id",
					Type:    ParamTypeString,
					Desc:    "UUID of Image",
					Default: ParamNotNull,
				},
			},
		},
	},
}
