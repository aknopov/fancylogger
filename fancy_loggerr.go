package fancylogger

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

//
// Based on https://github.com/rs/zerolog/issues/446
//

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite
	// colorBold     = 1
	// colorDarkGray = 90
)

const (
	timeFacet = "2006-01-02T15:04:05.000"
)

type CustomLogger struct {
	logger   zerolog.Logger
	curLevel any
}

var logger = NewLogger(os.Stdout)

func levelToColor(lvl any) int {
	switch lvl {
	case "info":
		return colorGreen
	case "warn":
		return colorYellow
	case "error", "fatal", "panic":
		return colorRed
	default:
		return colorBlack
	}
}

func colorize(s any, curLevel any) string {
	c := levelToColor(curLevel)
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

func colorizeFieldName(s any, curLevel any) string {
	return colorize(fmt.Sprintf("%s=", s), curLevel)
}

// Creates a new instance of custom logger.
// This instance shouild not be shared by go-routines
func NewLogger(writer io.Writer) CustomLogger {
	ret := CustomLogger{}

	colorizeLcl := func(s any) string {
		return colorize(s, ret.curLevel)
	}

	colorizeFieldLcl := func(s any) string {
		return colorizeFieldName(s, ret.curLevel)
	}

	customStandardOutput := zerolog.ConsoleWriter{
		Out:             writer,
		NoColor:         false,
		TimeFormat:      timeFacet,
		TimeLocation:    nil,
		PartsOrder:      []string{"time", "level", "application", "function", "message"},
		PartsExclude:    nil,
		FieldsOrder:     nil,
		FieldsExclude:   []string{"application", "function"},
		FormatTimestamp: nil,
		FormatLevel: func(i any) string {
			ret.curLevel = i
			return colorizeLcl(strings.ToUpper(fmt.Sprintf("%-5s|", i)))
		},
		FormatCaller:        nil,
		FormatMessage:       colorizeLcl,
		FormatFieldName:     colorizeFieldLcl,
		FormatFieldValue:    colorizeLcl,
		FormatErrFieldName:  colorizeFieldLcl,
		FormatErrFieldValue: colorizeLcl,
		FormatExtra:         nil,
		FormatPrepare:       nil,
	}

	zerolog.TimeFieldFormat = timeFacet

	ret.logger = zerolog.New(customStandardOutput).With().Timestamp().
		Str("application", "").
		Str("function", "").
		Logger()

	return ret
}

// Convenience adapters

func (l *CustomLogger) Trace() *zerolog.Event {
	return l.logger.Trace()
}

func (l *CustomLogger) Debug() *zerolog.Event {
	return l.logger.Debug()
}

func (l *CustomLogger) Info() *zerolog.Event {
	return l.logger.Info()
}

func (l *CustomLogger) Warn() *zerolog.Event {
	return l.logger.Warn()
}

func (l *CustomLogger) Error() *zerolog.Event {
	return l.logger.Error()
}

func (l *CustomLogger) Err(err error) *zerolog.Event {
	return l.logger.Err(err)
}

func (l *CustomLogger) Fatal() *zerolog.Event {
	return l.logger.Fatal()
}

func (l *CustomLogger) Panic() *zerolog.Event {
	return l.logger.Panic()
}
