package fancylogger

import (
	"errors"
	"math/rand"
	"regexp"
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
	tsRex, _ = regexp.Compile(TS_REGEX)
)

func TestSampleOutput(t *testing.T) {
	logger.Info().Msg("Hello log!")
}

func TestLevelToColor(t *testing.T) {
	assertT := assert.New(t)

	assertT.Equal(colorBlack, levelToColor("debug"))
	assertT.Equal(colorRed, levelToColor("error"))
	assertT.Equal(colorRed, levelToColor("panic"))
	assertT.Equal(colorGreen, levelToColor("info"))
	assertT.Equal(colorYellow, levelToColor("warn"))
	assertT.Equal(colorBlack, levelToColor("whatever"))
}

func TestColorize(t *testing.T) {
	assertColor(t, blackMarker, "debug")
	assertColor(t, greenMarker, "info")
	assertColor(t, yellowMarker, "warn")
	assertColor(t, redMarker, "error")
	assertColor(t, redMarker, "panic")
	assertColor(t, redMarker, "fatal")
}

func assertColor(t *testing.T, colorMarker []byte, lvl string) {
	assertT := assert.New(t)

	byteArr := []byte(colorize("", lvl))

	assertT.Equal(colorMarker, byteArr[:5])
	assertT.Equal(resetMarker, byteArr[5:])
}

func TestPlainLogging(t *testing.T) {
	assertT := assert.New(t)

	buffer := new(mockBuffer)
	testLogger := NewLogger(buffer)

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
}

func TestErrorLogging(t *testing.T) {
	assertT := assert.New(t)

	buffer := new(mockBuffer)
	testLogger := NewLogger(buffer)

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

func TestAdapters(t *testing.T) {
	assertT := assert.New(t)

	buffer := new(mockBuffer)
	testLogger := NewLogger(buffer)

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
	testLogger := NewLogger(buffer)
	rex, _ := regexp.Compile(TS_REGEX)

	for i := 0; i < 200; i++ {
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		testLogger.Info().Msg("")
		assertT.True(rex.MatchString(buffer.msg))
	}
}
