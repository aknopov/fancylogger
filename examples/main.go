package main

import (
	"os"

	fancylogger "github.com/aknopov/fancylogger/pkg"
)

func main() {
	lightLogger := fancylogger.NewLogger(os.Stdout, fancylogger.LiteFg)
	lightLogger.Info().Msg("Hello log!")

	darkLogger := fancylogger.NewLogger(os.Stdout, fancylogger.DarkFg)
	darkLogger.Warn().Msg("Hello warning!")

	noColorLogger := fancylogger.NewLogger(os.Stdout, fancylogger.NoColor)
	noColorLogger.Warn().Msg("Hello plain log!")
}
