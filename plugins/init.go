package plugins

import "octlink/ovs/utils/octlog"

var logger *octlog.LogConfig

// InitLog to init log config
func InitLog(level int) {
	logger = octlog.InitLogConfig("plugins.log", level)
}
