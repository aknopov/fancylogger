# Simple logger

Sample use -
```
    logger := NewLogger(os.Stdout, true)
    logger.Info().Msg("Hello log!")
```
The output -
<span style="font-weight: bold"><span style="color:gray">2025-04-05T13:20:02.059 </span><span style="color:green">INFO | Hello log!</span></span>

To prevent insering ANSI terminal sequences in the output, use `false` as the second parameter in the log creation.

 ## Notes
- Provides colorized entries on the console.
- Uses ISO timestamp with milliseconds.
- NOT thread-safe - each go-routine should create own instance.