package log

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/cloud-disk/internal/config"
	"github.com/cloud-disk/internal/constants"
)

type ZapLog struct {
	level    zap.AtomicLevel
	sugarLog *zap.SugaredLogger
	config   *config.LogConfig
}

func NewZapLog(conf *config.LogConfig) *ZapLog {
	zapLog := new(ZapLog)
	zapLog.config = conf
	zapLog.level = zap.NewAtomicLevel()
	zapLog.level.SetLevel(mapToLoggerLevel(conf.Level))

	writeSyncer := getLogWriter(conf)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapLog.level)

	zapLog.sugarLog = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
	return zapLog
}

func mapToLoggerLevel(level string) zapcore.Level {
	logLevelMap := map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
		"fatal": zapcore.FatalLevel,
	}

	level = strings.ToLower(level)
	ret, isExist := logLevelMap[level]
	if isExist {
		return ret
	}
	return zapcore.InfoLevel
}

func getLogWriter(conf *config.LogConfig) zapcore.WriteSyncer {
	syncers := []zapcore.WriteSyncer{
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(constants.RootPath, conf.ErrLogPath), // ⽇志⽂件路径
			MaxSize:    conf.MaxSize,                                       // 单位为MB,默认为512MB
			MaxAge:     conf.MaxAge,                                        // 文件最多保存多少天
			Compress:   conf.Compress,                                      // 是否压缩日志
			MaxBackups: conf.MaxBackups,                                    // 保存旧日志的文件数量
			LocalTime:  true,                                               // 是否使用当地时间
		}),
	}

	if conf.StdOut {
		syncers = append(syncers, zapcore.AddSync(os.Stdout))
	}
	return zapcore.NewMultiWriteSyncer(syncers...)
}

func getEncoder() zapcore.Encoder {
	// 自定义时间输出格式
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(constants.TimeFormat))
	}

	// 自定义日志级别显示
	customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(level.CapitalString())
	}

	encoderConf := zapcore.EncoderConfig{
		CallerKey:        "Caller",
		LevelKey:         "Level",
		MessageKey:       "Message",
		TimeKey:          "Time",
		StacktraceKey:    "Stack",
		NameKey:          "Name",
		LineEnding:       zapcore.DefaultLineEnding,
		ConsoleSeparator: "|",
		EncodeTime:       customTimeEncoder,
		EncodeLevel:      customLevelEncoder,
		EncodeDuration:   zapcore.StringDurationEncoder,
		EncodeName:       zapcore.FullNameEncoder,
	}

	return zapcore.NewConsoleEncoder(encoderConf)
}

func (z *ZapLog) close() {
	if z.sugarLog == nil {
		return
	}

	if err := z.sugarLog.Sync(); err != nil {
		fmt.Println(err)
	}
}

func appendPrefix(logFormat string) string {
	function, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}

	file = path.Base(file)
	funcName := runtime.FuncForPC(function).Name()

	// 每条日志信息格式为"文件名:行数|函数名|日志"
	builder := strings.Builder{}
	builder.WriteString(file)
	builder.WriteString(":")
	builder.WriteString(strconv.Itoa(line))
	builder.WriteString("|")
	builder.WriteString(funcName)
	builder.WriteString("|")
	builder.WriteString(logFormat)

	return builder.String()
}
