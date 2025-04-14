package fancylogger

import (
	"fmt"
	"io"
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

type ColorScheme int

// Types of color scheme.
const (
	NoColor ColorScheme = iota
	LiteFg
	DarkFg
)

func levelToColor(lvl any, colorScheme ColorScheme) int {
	switch lvl {
	case "info":
		return colorGreen
	case "warn":
		return colorYellow
	case "error", "fatal", "panic":
		return colorRed
	default:
		if colorScheme == LiteFg {
			return colorWhite
		}
		return colorBlack
	}
}

func colorize(s any, curLevel any, colorScheme ColorScheme) string {
	if s != nil {
		if colorScheme != NoColor {
			c := levelToColor(curLevel, colorScheme)
			return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
		}
		return fmt.Sprintf("%v", s)
	}
	return ""
}

func colorizeFieldName(s any, curLevel any, colorScheme ColorScheme) string {
	text := fmt.Sprintf("%s=", s)
	if colorScheme != NoColor {
		return colorize(text, curLevel, colorScheme)
	}
	return text
}

// Creates a new instance of custom logger.
//
// Color scheme sets color of non-colored text:
// 'LiteFg' uses teminal "white" color; 'DarkFg' uses "black" color
//
// This instance should not be shared by go-routines
func NewLogger(writer io.Writer, colorScheme ColorScheme) CustomLogger {
	ret := CustomLogger{}

	colorizeLcl := func(s any) string {
		return colorize(s, ret.curLevel, colorScheme)
	}

	colorizeFieldLcl := func(s any) string {
		return colorizeFieldName(s, ret.curLevel, colorScheme)
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
		FormatTimestamp: func(i any) string { return colorize(i, "", colorScheme) },
		FormatLevel: func(i any) string {
			ret.curLevel = i
			return colorizeLcl(strings.ToUpper(fmt.Sprintf("%-5s|", i)))
		},
		FormatCaller:        colorizeLcl,
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
