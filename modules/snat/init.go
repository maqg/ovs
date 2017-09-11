package snat

import "octlink/ovs/utils/octlog"

const (
	// MaxSnatsCount for max images count
	MaxSnatsCount = 1000
)

// GSnats for all image loaded from config
var GSnats []*Snat

func loadSnatsFromConfig() error {
	GSnats = make([]*Snat, 0)
	return nil
}

func zeroSnats() {
	GSnats = make([]*Snat, 0)
}

// ReloadImages for images reloading
func ReloadImages() error {

	// zero images firstly
	zeroSnats()

	err := loadSnatsFromConfig()
	if err != nil {
		octlog.Error("load images error [%s]\n", err)
		return nil
	}

	return nil
}

// WriteImages to write all images to image store file
func WriteImages() error {
	// TBD
	return nil
}
