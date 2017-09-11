package snat

import "octlink/ovs/utils/octlog"

var logger *octlog.LogConfig

// InitLog to init log config
func InitLog(level int) {
	logger = octlog.InitLogConfig("image.log", level)
}

const (
	// SnatConfigFile for image basic info store file
	SnatConfigFile = "snat_config.json"
)

// Snat for snat sturcture
type Snat struct {
	ID            string `json:"id"`
	PrivateNicMac string `json:"privateNicMac"`
}

// GetSnatCount to return image count by condition
func GetSnatCount() int {
	return len(GetAllSnats())
}

// Brief to return brief info for image
func (s *Snat) Brief() map[string]string {
	return map[string]string{
		"id":   s.ID,
		"name": s.Name,
	}
}

// Update to update image
func (s *Snat) Update() int {
	WriteImages()
	return 0
}

// Add for image, after image added,
// installpath, diskSize, virtualSize, Status, md5sum need update after manifest installed
func (s *Snat) Add() int {
	return 0
}

// GetSnat get snat settings
func GetSnat(id string) *Snat {
	return nil
}

// GetAllSnats by condition
func GetAllSnats() []*Snat {
	return make([]*Snat, 0)
}
