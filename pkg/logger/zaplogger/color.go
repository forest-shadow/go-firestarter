package zaplogger

import (
	"strconv"

	"go.uber.org/zap/zapcore"
)

const (
	gruvboxGray   = 245
	gruvboxBlue   = 109
	gruvboxGreen  = 142
	gruvboxYellow = 214
	gruvboxOrange = 208
	gruvboxRed    = 167
	gruvboxReset  = "\x1b[0m"
)

func levelColor(level zapcore.Level) int {
	switch level {
	case zapcore.DebugLevel:
		return gruvboxBlue
	case zapcore.InfoLevel:
		return gruvboxGreen
	case zapcore.WarnLevel:
		return gruvboxYellow
	case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		return gruvboxRed
	default:
		return gruvboxOrange
	}
}

func applyColor(colorCode int, value string, enabled bool) string {
	if !enabled || value == "" {
		return value
	}

	return "\x1b[38;5;" + strconv.Itoa(colorCode) + "m" + value + gruvboxReset
}
