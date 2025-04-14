# Simple logger

- Provides colorized output on the console.
- Uses ISO timestamp with milliseconds.
- NOT thread-safe - each go-routine should create own instance.
- Color output can be switched off with `NoColor` option of log creation.
- Color output has two variants with light or dark foreground. Actual colors depend on terminal settings.

Sample use -
```
	lightLogger := NewLogger(os.Stdout, LiteFg)
	lightLogger.Info().Msg("Hello log!")

	darkLogger := NewLogger(os.Stdout, DarkFg)
	darkLogger.Info().Msg("Hello log!")
```
The output - <img style="vertical-align: top" src="./output_sample.png" alt="sample output"/>

 ## Notes
