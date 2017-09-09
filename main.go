package main

import (
	"flag"
	"fmt"
	"net/http"
	"octlink/ovs/api"
	"octlink/ovs/utils"
	"octlink/ovs/utils/configuration"
	"octlink/ovs/utils/octlog"
)

var (
	port   int
	addr   string
	config string
)

var conf *configuration.Configuration

func initDebugConfig() {
	octlog.InitDebugConfig(conf.DebugLevel)
}

func initLogConfig() {
	utils.CreateDir(conf.LogDirectory)
	api.InitLog(conf.LogLevel)
	utils.InitLog(conf.LogLevel)
}

func initDebugAndLog() {
	initDebugConfig()
	initLogConfig()
}

func init() {
	flag.StringVar(&config, "config", "./config.yml", "Config file path")
}

func usage() {
	fmt.Printf("  RVM Store of V" + utils.Version() + "\n")
	fmt.Printf("  ./ovs -config ./config.yml\n")
	flag.PrintDefaults()
}

func runAPIThread() {

	api := &api.API{
		Name: "OVS API Server",
	}

	server := &http.Server{
		Addr:           fmt.Sprintf("%s", conf.HTTP.Addr),
		Handler:        api.Router(),
		MaxHeaderBytes: 1 << 20,
	}

	octlog.Warn("OVS API Engine Started ON %s\n", conf.HTTP.Addr)

	err := server.ListenAndServe()
	if err != nil {
		octlog.Error("error to listen at %s\n", conf.HTTP.Addr)
	}
}

func main() {

	flag.Usage = usage
	flag.Parse()

	c, err := configuration.ResolveConfig(config)
	if err != nil {
		fmt.Printf("Resolve Configuration Error[%s]\n", err)
		return
	}
	conf = c

	initDebugAndLog()

	runAPIThread()
}
