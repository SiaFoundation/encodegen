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
encodegen -s . -t MyType -o output.go
```

### Using flags

    * -o: destination file to output generated code
    * -pkg: the package name of the generated file
    * -s: Source dir or file (absolute or relative path), omit for stdout
    * -t: Types to generate, comma separated.  To enable memory reuse, put "true" after a type, e.g. Message,true,SubMessage,SubMessage2.  Memory reuse defaults to false if not specified.

## Test

```
$ cd internal/codegen/test
$ ./generate_tests.sh
```

Tests for a basic struct (`basic_struct`), a struct using aliased types (`alias_struct`), and a struct with anonymous structs in it (`embedded_struct`) are now available.

```
$ cd basic_struct
$ go test
```