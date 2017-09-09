package snat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"octlink/ovs/utils"
	"octlink/ovs/utils/configuration"
	"octlink/ovs/utils/octlog"
	"os"
)

const (
	// MaxSnatsCount for max images count
	MaxSnatsCount = 1000
)

// GSnats for all image loaded from config
var GSnats []*Snat

func loadSnatsFromConfig() error {

	imagePath := configuration.GetConfig().RootDirectory + "/" + SnatConfigFile
	if !utils.IsFileExist(imagePath) {
		octlog.Error("file of %s not exist\n", imagePath)
		return fmt.Errorf("file of %s not exist", imagePath)
	}

	file, err := os.Open(imagePath)
	if err != nil {
		octlog.Error("open image store file " + imagePath + "error\n")
		return err
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	err = json.Unmarshal(data, &GSnats)
	if err != nil {
		octlog.Warn("Transfer json bytes error %s\n", err)
		return err
	}

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

	imagePath := configuration.GetConfig().RootDirectory + "/" + SnatConfigFile
	if utils.IsFileExist(imagePath) {
		os.Rename(imagePath, imagePath+"."+utils.CurrentTimeSimple())
	}

	fd, err := os.Create(imagePath)
	if err != nil {
		octlog.Error("create file of %s error\n", imagePath)
		return err
	}

	_, err = fd.Write(utils.JSON2Bytes(GSnats))
	if err != nil {
		octlog.Error("Write images to image store file %s error\n", imagePath)
		return err
	}

	return nil
}
