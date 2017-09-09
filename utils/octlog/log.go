package octlog

import (
	"fmt"
	"log"
	"octlink/ovs/utils/configuration"
	"os"
)

const (

	//PanicLevel for panic
	PanicLevel int = iota

	// FatalLevel for Fatal
	FatalLevel

	// ErrorLevel for Error
	ErrorLevel

	// WarnLevel for Warn
	WarnLevel

	// InfoLevel for Info
	InfoLevel

	// DebugLevel for Debug
	DebugLevel
)

// GDebugConfig for Global Debug Configuation
var GDebugConfig DebugConfig

// DebugConfig Structure
type DebugConfig struct {
	level int
}

// LogConfig Stucture
type LogConfig struct {
	level   int
	logTime int64
	LogFile string
	fileFd  *os.File
	logger  *log.Logger
}

func getLogDir() string {
	return configuration.LogDirectory()
}

// InitLogConfig to init log config
func InitLogConfig(logFile string, level int) *LogConfig {
	config := new(LogConfig)

	config.LogFile = getLogDir() + "/" + logFile
	config.level = level

	logfile, err := os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("%s\r\n", err.Error())
		return nil
	}

	logfile.Seek(0, 2)

	config.logger = log.New(logfile, "", log.Ldate|log.Ltime|log.Lshortfile)

	return config
}

// InitDebugConfig to init debug config
func InitDebugConfig(level int) {
	GDebugConfig.level = level
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
}

// Debugf For File Logging
func (config *LogConfig) Debugf(format string, args ...interface{}) {
	if config.level >= DebugLevel {
		config.logger.SetPrefix("DEBUG ")
		config.logger.Printf(format, args...)
	}
}

// Infof for Info File Logging
func (config *LogConfig) Infof(format string, args ...interface{}) {
	if config.level >= InfoLevel {
		config.logger.SetPrefix("INFO ")
		config.logger.Printf(format, args...)
	}
}

// Warnf for Warn File Logging
func (config *LogConfig) Warnf(format string, args ...interface{}) {
	if config.level >= WarnLevel {
		config.logger.SetPrefix("WARN ")
		config.logger.Printf(format, args...)
	}
}

// Errorf for Error File Logging
func (config *LogConfig) Errorf(format string, args ...interface{}) {
	if config.level >= ErrorLevel {
		config.logger.SetPrefix("ERROR ")
		config.logger.Printf(format, args...)
	}
}

// Fatalf for Fatal File Logging
func (config *LogConfig) Fatalf(format string, args ...interface{}) {
	if config.level >= FatalLevel {
		config.logger.SetPrefix("FATAL ")
		config.logger.Printf(format, args...)
	}
}

// Panicf for Panic File Logging
func (config *LogConfig) Panicf(format string, args ...interface{}) {
	if config.level >= PanicLevel {
		config.logger.SetPrefix("PANIC ")
		config.logger.Printf(format, args...)
	}
}

// Debug for Debuging
func Debug(format string, args ...interface{}) {
	if GDebugConfig.level >= DebugLevel {
		log.SetPrefix("DEBUG ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

// Info for infoing
func Info(format string, args ...interface{}) {
	if GDebugConfig.level >= InfoLevel {
		log.SetPrefix("INFO ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

// Warn for Warn Debug
func Warn(format string, args ...interface{}) {
	if GDebugConfig.level >= WarnLevel {
		log.SetPrefix("WARN ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

// Error for error debuging
func Error(format string, args ...interface{}) {
	if GDebugConfig.level >= ErrorLevel {
		log.SetPrefix("ERROR ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

// Fatal for fatal debuging
func Fatal(format string, args ...interface{}) {
	if GDebugConfig.level >= FatalLevel {
		log.SetPrefix("FATAL ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

// Panic for panic debuging
func Panic(format string, args ...interface{}) {
	if GDebugConfig.level >= PanicLevel {
		log.SetPrefix("PANIC ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}
