/**
 * @Author: jinpeng zhang
 * @Date: 2023/8/13 14:29
 * @Description:
 */

package logger

type Config struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
	Name              string
	Writers           string
	LoggerFile        string
	LoggerWarnFile    string
	LoggerErrorFile   string
	LogFormatText     bool
	LogRollingPolicy  string
	LogBackupCount    uint
}
