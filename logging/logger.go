package logging

import (
	"github.com/gin-gonic/gin"
	golog "github.com/umisama/golog"
	"os"
)

const logFormat = "{{.Time}} {{.FileName}}:{{.LineNumber}}({{.FuncName}}) : {{.Message}}\n"

var logger golog.Logger

// Logger returns a logger instance
func Logger() golog.Logger {
	return logger
}

func init() {
	if gin.IsDebugging() {
		logger, _ = golog.NewLogger(os.Stdout,
			golog.TIME_FORMAT_SEC,
			logFormat,
			golog.LogLevel_Debug)
	} else {
		logger, _ = golog.NewLogger(os.Stdout,
			golog.TIME_FORMAT_SEC,
			logFormat,
			golog.LogLevel_Critical)
	}
}
