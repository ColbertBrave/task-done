package log

import (
	"path/filepath"

	"cloud-disk/internal/config"
	"cloud-disk/internal/constants"
)

var (
	sysLog *SystemLog
	zapLog *ZapLog
)

func InitLog(conf *config.LogConfig) {
	sysLog = NewSystemLog(filepath.Join(constants.RootPath, config.AppCfg.LogCfg.SysLogPath))
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

	prefix := appendPrefix()
	if len(args) == 0 {
		zapLog.sugarLog.Error(prefix + format)
		return
	}

	zapLog.sugarLog.Errorf(prefix+format, args)
}

func Fatal(format string, args ...interface{}) {
	if zapLog == nil {
		return
	}

	prefix := appendPrefix()
	if len(args) == 0 {
		zapLog.sugarLog.Fatal(prefix + format)
		return
	}

	zapLog.sugarLog.Fatalf(prefix+format, args)
}

func Info(format string, args ...interface{}) {
	if zapLog == nil {
		return
	}

	prefix := appendPrefix()
	if len(args) == 0 {
		zapLog.sugarLog.Infof(prefix + format)
		return
	}

	zapLog.sugarLog.Infof(prefix+format, args)
}

func Debug(format string, args ...interface{}) {
	if zapLog == nil {
		return
	}

	prefix := appendPrefix()
	if len(args) == 0 {
		zapLog.sugarLog.Debug(prefix + format)
		return
	}

	zapLog.sugarLog.Debugf(prefix+format, args)
}

func Warn(format string, args ...interface{}) {
	if zapLog == nil {
		return
	}

	prefix := appendPrefix()
	if len(args) == 0 {
		zapLog.sugarLog.Warn(prefix + format)
		return
	}

	zapLog.sugarLog.Warnf(prefix+format, args)
}
