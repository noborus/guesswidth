# guesswidth

[![Go](https://github.com/noborus/guesswidth/actions/workflows/build.yml/badge.svg)](https://github.com/noborus/guesswidth/actions/workflows/build.yml)

## ovewerview

guesswidth guesses the separator position in CLI output.

The output of the `ps` command has no delimiters, making the values difficult to machine-readable.
guesswidth is smarter at guessing separators than just spaces.

## install command

There is also a guesswidth command.

```console
go install github.com/noborus/guesswidth/cmd/guesswidth
```

Guess the width output with no delimiters in the command.

```console
$ ps
    PID TTY          TIME CMD
 302965 pts/3    00:00:12 zsh
 733211 pts/3    00:00:00 ps
 733212 pts/3    00:00:00 tee
 733213 pts/3    00:00:00 guesswidth
```

Split the output like this:

```console
$ ps | guesswidth
    PID| TTY     |     TIME|CMD
 302965| pts/3   | 00:00:08|zsh
 539529| pts/3   | 00:00:00|ps
 539530| pts/3   | 00:00:00|guesswidth
```

It can be converted to `csv`.

```console
$ ps | guesswidth csv
PID,TTY,TIME,CMD
302965,pts/3,00:00:12,zsh
733211,pts/3,00:00:00,ps
733212,pts/3,00:00:00,tee
733213,pts/3,00:00:00,guesswidth
```
