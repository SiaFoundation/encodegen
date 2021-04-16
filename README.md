# encodegen

`encodegen` generates marshaling and unmarshaling implementations for the
[Sia encoding format](https://gitlab.com/NebulousLabs/encoding).

## Installation

```
go install go.sia.tech/encodegen
```

## Usage

Add a `//go:generate` directive to the package containing the types you want to
generate methods for:

```go
//go:generate encodegen -t Foo,Bar

type Foo struct {
    ...
}

type Bar struct {
    ...
}
```

When run, this will add an `encoding.go` file to the package, containing
`MarshalSia` and `UnmarshalSia` method definitions for each of the specified
types.
