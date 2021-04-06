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