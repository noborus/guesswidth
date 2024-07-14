# guesswidth

[![Go Reference](https://pkg.go.dev/badge/github.com/noborus/guesswidth.svg)](https://pkg.go.dev/github.com/noborus/guesswidth) [![Go](https://github.com/noborus/guesswidth/actions/workflows/build.yml/badge.svg)](https://github.com/noborus/guesswidth/actions/workflows/build.yml)

## Overview

guesswidth guesses the column position for fixed-width formats(fwf).

The output of the `ps` command has no delimiters, making the values difficult to machine-readable.
guesswidth guesses smarter than just space delimiters.

guesswidth is guessed based on the position of characters in the **header**.
So having a header will give you better results.

## Install command

There is also a guesswidth command.

```console
go install github.com/noborus/guesswidth/cmd/guesswidth@latest
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

```console
ps | guesswidth --align
   PID|TTY  |    TIME|CMD
302965|pts/3|00:00:12|zsh
733211|pts/3|00:00:00|ps
733212|pts/3|00:00:00|tee
733213|pts/3|00:00:00|guesswidth
```

## Usage

```console
$ guesswidth -h         [~/dev/src/github.com/noborus/guesswidth]
Guess the width of the columns from the header and body,
split them, insert fences and output.

Usage:
  guesswidth [flags]
  guesswidth [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  csv         Output in csv format
  help        Help about any command

Flags:
  -a, --align          align the output
      --fence string   fence (default "|")
      --header int     header line number (default 1)
  -h, --help           help for guesswidth
      --scannum int    number of line to scan (default 100)
      --split int      maximum number of splits (default -1)
  -v, --version        version for guesswidth

Use "guesswidth [command] --help" for more information about a command.
```

## Examples

guesswidth inserts a delimiter (| by default)
(Colors are changed here for clarity).

Even if there are spaces in the header or body, they will be separated correctly.

### ps

![ps](./docs/ps.png)

### docker ps

![docker-ps](./docs/docker-ps.png)

### docker node

![docker node ls](./docs/docker-node.png)