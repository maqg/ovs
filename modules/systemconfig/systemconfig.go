package systemconfig

import (
	"octlink/ovs/utils/configuration"
)

// SystemConfig for system
type SystemConfig struct {
	Version string `json:"version"`
	Snat    int    `json:"snat"`
	Vip     int    `json:"vip"`
	Eip     int    `json:"eip"`
}

// GetSystemConfig get system config of this backupstorage
func GetSystemConfig() *SystemConfig {

	conf := configuration.GetConfig()

	sc := new(SystemConfig)
	sc.Version = conf.Version
	sc.Eip = 10
	sc.Vip = 100
	sc.Snat = 320

	return sc
}
