package api

var nicDescriptors = Module{
	Name: "nic",
	Protos: map[string]Proto{
		"APIShowInterfaces": {
			Name:    "查看接口信息",
			handler: ShowInterfaces
			Paras:   []ProtoPara{},
		},
	},
}
