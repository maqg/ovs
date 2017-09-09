package api

// ConfigDescriptors for image management by API
var ConfigDescriptors = Module{
	Name: "config",
	Protos: map[string]Proto{
		"APIShowSystemInfo": {
			Name:    "查看系统信息",
			handler: ShowSystemConfig,
			Paras:   []ProtoPara{},
		},
	},
}
