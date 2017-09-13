package plugins

import "octlink/ovs/utils/octlog"

var logger *octlog.LogConfig

const (
	// EipSnatStartRuleNum for snat rule
	EipSnatStartRuleNum = 5000
)

// InitLog to init log config
func InitLog(level int) {
	logger = octlog.InitLogConfig("plugins.log", level)
}
