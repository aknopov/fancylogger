package fancylogger

import (
	"errors"
	"math/rand"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockBuffer struct {
	msg string
}

func (b *mockBuffer) Write(p []byte) (n int, err error) {
	b.msg = string(p)
	return len(p), nil
}

const (
	TS_REGEX = "\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}\\.\\d{3}"
)

var (
	resetMarker  = []byte{27, 91, 48, 109}
	blackMarker  = []byte{27, 91, 51, 48, 109}
	greenMarker  = []byte{27, 91, 51, 50, 109}
	yellowMarker = []byte{27, 91, 51, 51, 109}
	redMarker    = []byte{27, 91, 51, 49, 109}
	whiteMarker  = []byte{27, 91, 51, 55, 109}
	tsRex        = regexp.MustCompile(TS_REGEX)
)

func TestLevelToColor(t *testing.T) {
	assertT := assert.New(t)

	tests := []struct {
		level    string
		scheme   ColorScheme
		colorIdx int
	}{
		{"debug", DarkFg, colorBlack},
		{"debug", LiteFg, colorWhite},
		{"error", DarkFg, colorRed},
		{"error", LiteFg, colorRed},
		{"panic", DarkFg, colorRed},
		{"panic", LiteFg, colorRed},
		{"info", DarkFg, colorGreen},
		{"info", LiteFg, colorGreen},
		{"warn", DarkFg, colorYellow},
		{"warn", LiteFg, colorYellow},
		{"whatever", DarkFg, colorBlack},
		{"whatever", LiteFg, colorWhite},
	}

	for _, tt := range tests {
		assertT.Equal(tt.colorIdx, levelToColor(tt.level, tt.scheme))
	}
}

func TestColorize(t *testing.T) {
	assertT := assert.New(t)

	tests := []struct {
		level       string
		scheme      ColorScheme
		colorMarker []byte
	}{
		{"debug", DarkFg, blackMarker},
		{"debug", LiteFg, whiteMarker},
		{"info", DarkFg, greenMarker},
		{"info", LiteFg, greenMarker},
		{"warn", DarkFg, yellowMarker},
		{"warn", LiteFg, yellowMarker},
		{"error", DarkFg, redMarker},
		{"error", LiteFg, redMarker},
		{"panic", DarkFg, redMarker},
		{"panic", LiteFg, redMarker},
		{"fatal", DarkFg, redMarker},
		{"fatal", LiteFg, redMarker},
	}

	for _, tt := range tests {
		byteArr := []byte(colorize("a", tt.level, tt.scheme))

		assertT.Equal(tt.colorMarker, byteArr[:5])
		assertT.Equal(resetMarker, byteArr[6:])
	}
}

func TestNocolorForEmptyString(t *testing.T) {
	assertT := assert.New(t)

	tests := []struct {
		val any
		fg  ColorScheme
	}{
		{nil, DarkFg},
		{nil, LiteFg},
		{nil, NoColor},
		{"", DarkFg},
		{"", LiteFg},
		{"", NoColor},
	}
	for _, tt := range tests {
		assertT.Equal("", colorize(tt.val, "info", tt.fg))
	}
}

func TestNoErrorLogging(t *testing.T) {
	assertT := assert.New(t)

	buffer := new(mockBuffer)
	testLogger := NewLogger(buffer, LiteFg)

	testLogger.logger.Info().
		Str("Param", "String value").
		Msg("Here you are:")

	logEntry := buffer.msg
	assertT.Subset([]byte(logEntry), greenMarker)
	assertT.True(tsRex.MatchString(logEntry))
	assertT.Contains(logEntry, "INFO |")
	assertT.Contains(logEntry, "Here you are:")
	assertT.Contains(logEntry, "Param=")
	assertT.Contains(logEntry, "String value")
	assertT.True(strings.HasSuffix(logEntry, "\n"))
}

func TestErrorLogging(t *testing.T) {
	assertT := assert.New(t)

	buffer := new(mockBuffer)
	testLogger := NewLogger(buffer, LiteFg)

	testLogger.Error().
		Err(errors.New("NFG")).
		Msg("Here you are:")

	logEntry := buffer.msg
	assertT.Subset([]byte(logEntry), redMarker)
	assertT.True(tsRex.MatchString(logEntry))
	assertT.Contains(logEntry, "ERROR|")
	assertT.Contains(logEntry, "Here you are:")
	assertT.Contains(logEntry, "error=")
	assertT.Contains(logEntry, "NFG")
}

func TestNoColor(t *testing.T) {
	assertT := assert.New(t)

	buffer := new(mockBuffer)
	testLogger := NewLogger(buffer, NoColor)
	testLogger.logger.Info().
		Str("Param", "String value").
		Msg("Here you are:")

	logEntry := buffer.msg
	assertT.NotContains(logEntry, "\x1b[")
	assertT.True(tsRex.MatchString(logEntry))
	assertT.Contains(logEntry, "Here you are:")
	assertT.Contains(logEntry, "Param=")
	assertT.Contains(logEntry, "String value")
	assertT.True(strings.HasSuffix(logEntry, "\n"))
}

func TestAdapters(t *testing.T) {
	assertT := assert.New(t)

	buffer := new(mockBuffer)
	testLogger := NewLogger(buffer, LiteFg)

	testLogger.Trace().Msg("")
	assertT.Contains(buffer.msg, "TRACE")
	testLogger.Debug().Msg("")
	assertT.Contains(buffer.msg, "DEBUG")
	testLogger.Info().Msg("")
	assertT.Contains(buffer.msg, "INFO")
	testLogger.Warn().Msg("")
	assertT.Contains(buffer.msg, "WARN")
	testLogger.Error().Msg("")
	assertT.Contains(buffer.msg, "ERROR")
	testLogger.Err(errors.New("NFG"))
	assertT.Contains(buffer.msg, "ERROR")
	assertT.Panics(func() { testLogger.Panic().Msg("") })
	assertT.Contains(buffer.msg, "PANIC")
}

func TestTimestamp(t *testing.T) {
	assertT := assert.New(t)

	buffer := new(mockBuffer)
	testLogger := NewLogger(buffer, LiteFg)

	for i := 0; i < 200; i++ {
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		testLogger.Info().Msg("")
		assertT.True(tsRex.MatchString(buffer.msg))
	}
}

func TestNoNilMessage(t *testing.T) {
	assertT := assert.New(t)

	buffer := new(mockBuffer)
	testLogger := NewLogger(buffer, LiteFg)

	testLogger.Info().Dur("param", time.Duration(1234567)).Send()
	logEntry := buffer.msg
	assertT.NotContains(logEntry, "<nil>")
	assertT.Contains(logEntry, "INFO |")
	assertT.Contains(logEntry, "param")
	assertT.Contains(logEntry, "1.234567")
}
