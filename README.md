# Simple logger

Sample use -
```
    logger := NewLogger(os.Stdout)
    logger.Info().Msg("Hello log!")
```
The output -
<span style="font-weight: bold"><span style="color:gray">2025-04-05T13:20:02.059 </span><span style="color:green">INFO |&nbsp;&nbsp;&nbsp;Hello log!</span></span>

 ## Notes
- Provides colorized entries on the console.
- Uses ISO timestamp with milliseconds.
- NOT thread-safe - each go-routine should create own instance.