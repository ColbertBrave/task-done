package log

import (
	"fmt"
	"path/filepath"

	"github.com/cloud-disk/infrastructure/config"
	"github.com/cloud-disk/infrastructure/constants"
)

var (
	sysLog *SystemLog
	zapLog *ZapLog
)

func InitLog(conf *config.LogConfig) {
	sysLog = NewSystemLog(filepath.Join(constants.RootPath, config.GetConfig().LogCfg.SysLogPath))
	zapLog = NewZapLog(conf)
}

func Close() {
	if sysLog != nil {
		sysLog.close()
		sysLog = nil
	}

	if zapLog != nil {
		zapLog.close()
		zapLog = nil
	}
}

func Error(format string, args ...interface{}) {
	if zapLog == nil {
		return
	}

	prefix := appendPrefix(format)
	if len(args) == 0 {
		zapLog.sugarLog.Error(prefix)
		return
	}

	zapLog.sugarLog.Errorf(prefix, args)
}

func Fatal(format string, args ...interface{}) {
	if zapLog == nil {
		return
	}

	prefix := appendPrefix(format)
	if len(args) == 0 {
		zapLog.sugarLog.Fatal(prefix)
		return
	}

	zapLog.sugarLog.Fatalf(prefix, args)
}

func Info(format string, args ...interface{}) {
	if zapLog == nil {
		return
	}

	prefix := appendPrefix(format)
	if len(args) == 0 {
		zapLog.sugarLog.Info(prefix)
		return
	}

	zapLog.sugarLog.Infof(prefix, args)
}

func Debug(format string, args ...interface{}) {
	if zapLog == nil {
		return
	}

	prefix := appendPrefix(format)
	if len(args) == 0 {
		zapLog.sugarLog.Debug(prefix)
		return
	}

	zapLog.sugarLog.Debugf(prefix, args)
}

func Warn(format string, args ...interface{}) {
	if zapLog == nil {
		return
	}

	prefix := appendPrefix(format)
	if len(args) == 0 {
		zapLog.sugarLog.Warn(prefix)
		return
	}

	zapLog.sugarLog.Warnf(prefix, args)
}

func Sys(format string, args ...interface{}) {
	if sysLog == nil {
		fmt.Printf(format, args...)
		return
	}

	sysLog.log(format, args)
}
