# guesswidth

## install

```console
go install github.com/noborus/guesswidth/cmd/guesswidth
```

Guess the width output with no delimiters in the command.

Split the output like this:

```console
$ ps | guesswidth
    PID| TTY     |     TIME|CMD
 302965| pts/3   | 00:00:08|zsh
 539529| pts/3   | 00:00:00|ps
 539530| pts/3   | 00:00:00|guesswidth
```
