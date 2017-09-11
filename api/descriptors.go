package api

var apiDescriptors = []Module{
	snatDescriptors,
	configDescriptors,
}

func loadModules(module Module) {

	if GAPIConfig.Modules == nil {
		GAPIConfig.Modules = make(map[string]Module, 50)
	}

	GAPIConfig.Modules[module.Name] = module
}

func init() {

	GServices = make(map[string]*Service, 10000)

	for _, descriptor := range apiDescriptors {
		for key, proto := range descriptor.Protos {
			service := new(Service)
			service.Name = proto.Name
			service.Handler = proto.handler
			proto.Key = APIPrefixCenter + "." + descriptor.Name + "." + key
			descriptor.Protos[key] = proto
			GServices[APIPrefixCenter+"."+descriptor.Name+"."+key] = service
		}
		loadModules(descriptor)
	}
}
