/**
 * @Author: jinpeng zhang
 * @Date: 2023/8/13 14:30
 * @Description:
 */

package logger

import (
	"blog/pkg/utils"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"strings"
	"time"
)

const (
	WriterConsole = "console"
	WriterField   = "file"
)

const (
	RotateTimeDaily  = "daily"
	RotateTimeHourly = "hourly"
)

var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func getLoggerLevel(cfg *Config) zapcore.Level {
	level, exist := loggerLevelMap[cfg.Level]
	if !exist {
		return zapcore.DebugLevel
	}
	return level
}

func newZapLogger(cfg *Config) (*zap.Logger, error) {
	return buildLogger(cfg), nil
}

func buildLogger(cfg *Config) *zap.Logger {
	var encoderCfg zapcore.EncoderConfig
	if cfg.Development {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}
	// 自定义时间输出格式
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2020-01-01 00:00:00.000"))
	}
	encoderCfg.EncodeTime = customTimeEncoder
	var encoder zapcore.Encoder
	if cfg.Encoding == WriterConsole {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}
	var cores []zapcore.Core
	var options []zap.Option
	hostname, _ := os.Hostname()
	option := zap.Fields(
		zap.String("ip", utils.GetLocalIp()),
		zap.String("app_id", cfg.Name),
		zap.String("instance_id", hostname),
	)
	options = append(options, option)
	writers := strings.Split(cfg.Writers, ",")
	for _, w := range writers {
		switch w {
		case WriterConsole:
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), getLoggerLevel(cfg)))
		case WriterField:
			// info
			cores = append(cores, getInfoCore(encoder, cfg))
			// warning
			core, option := getWarnCore(encoder, cfg)
			cores = append(cores, core)
			if utils.IsNotNil(option) {
				options = append(options, option)
			}
			// error
			core, option = getErrorCore(encoder, cfg)
			cores = append(cores, core)
			if utils.IsNotNil(option) {
				options = append(options, option)
			}
		default:
			// console
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), getLoggerLevel(cfg)))
			// file
			cores = append(cores, getAllCore(encoder, cfg))
		}
	}
	combinedCore := zapcore.NewTee(cores...)
	// 开启开发模式，堆栈跟踪
	if !cfg.DisableCaller {
		caller := zap.AddCaller()
		options = append(options, caller)
	}
	// 跳过文件调用层数
	addCallerSkip := zap.AddCallerSkip(2)
	options = append(options, addCallerSkip)
	return zap.New(combinedCore, options...)
}

func getAllCore(encoder zapcore.Encoder, cfg *Config) zapcore.Core {
	allWriter := getLogWriterWithTime(cfg, cfg.LoggerFile)
	allLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.FatalLevel
	})
	return zapcore.NewCore(encoder, zapcore.AddSync(allWriter), allLevel)
}

func getInfoCore(encoder zapcore.Encoder, cfg *Config) zapcore.Core {
	infoWrite := getLogWriterWithTime(cfg, cfg.LoggerFile)
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.InfoLevel
	})
	return zapcore.NewCore(encoder, zapcore.AddSync(infoWrite), infoLevel)
}

func getWarnCore(encoder zapcore.Encoder, cfg *Config) (zapcore.Core, zap.Option) {
	warnWrite := getLogWriterWithTime(cfg, cfg.LoggerWarnFile)
	var stacktrace zap.Option
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if !cfg.DisableCaller {
			stacktrace = zap.AddStacktrace(zapcore.WarnLevel)
		}
		return lvl == zapcore.WarnLevel
	})
	return zapcore.NewCore(encoder, zapcore.AddSync(warnWrite), warnLevel), stacktrace
}

func getErrorCore(encoder zapcore.Encoder, cfg *Config) (zapcore.Core, zap.Option) {
	errorFilename := cfg.LoggerErrorFile
	errorWrite := getLogWriterWithTime(cfg, errorFilename)
	var stacktrace zap.Option
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if !cfg.DisableCaller {
			stacktrace = zap.AddStacktrace(zapcore.ErrorLevel)
		}
		return lvl >= zapcore.ErrorLevel
	})
	return zapcore.NewCore(encoder, zapcore.AddSync(errorWrite), errorLevel), stacktrace
}

func getLogWriterWithTime(cfg *Config, filename string) io.Writer {
	logFullPath := filename
	rotationPolicy := cfg.LogRollingPolicy
	backupCount := cfg.LogBackupCount
	var rotateDuration time.Duration
	if rotationPolicy == RotateTimeHourly {
		rotateDuration = time.Hour
	} else {
		rotateDuration = time.Hour * 24
	}
	hook, err := rotatelogs.New(
		logFullPath+".%Y%m%d%H",
		rotatelogs.WithLinkName(logFullPath),
		rotatelogs.WithRotationCount(backupCount),
		rotatelogs.WithRotationTime(rotateDuration),
	)
	if utils.IsNotNil(err) {
		panic(err)
	}
	return hook
}
