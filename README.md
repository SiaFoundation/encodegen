# Encodegen
This package provides a command line tool to generate [NebulousLabs/encoding](https://gitlab.com/NebulousLabs/encoding) marshaling and unmarshaling interface implementations for custom struct types.

## Get started

```
go install go.sia.tech/encodegen
```

The `encodegen` binary will now appear in your `$GOPATH/bin` directory.

## Generate code

### Basic command
The basic command is easy to use:
```
encodegen -t MyType -o output.go
```

### Flags

    * -pkg: name of target package
    * -o: destination of generated code (optional; omit for stdout)
    * -t: types to generate, comma separated

## Test

    $ cd test
    $ go run ../. -t TestMessageSimple,TestMessageEmbedded -o encoding.go
    $ go test

## Fuzzing

Ensure [go-fuzz](https://github.com/dvyukov/go-fuzz) is installed:

    $ go get -u github.com/dvyukov/go-fuzz/go-fuzz github.com/dvyukov/go-fuzz/go-fuzz-build


Generate corpus:

    $ cd test/fuzz/gencorpus
    $ go run . -out ../corpus
    $ cd ..

Instrument the code:

    $ go-fuzz-build

Begin fuzzing

    $ go-fuzz

You should now see output like

    2021/04/14 10:18:14 workers: 4, corpus: 100 (33s ago), crashers: 0, restarts: 1/9589, execs: 795925 (24105/sec), cover: 0, uptime: 33s
