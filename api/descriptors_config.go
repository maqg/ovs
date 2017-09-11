package api

// configDescriptors for image management by API
var configDescriptors = Module{
	Name: "config",
	Protos: map[string]Proto{
		"APIShowSystemInfo": {
			Name:    "查看系统信息",
			handler: ShowSystemConfig,
			Paras:   []ProtoPara{},
		},
	},
}
