package codegen

const ReadBoolFunction = "ReadBool"
const ReadIntFunction = "ReadUint64"
const ReadStringFunction = "ReadString"
const ReadByteFunction = "ReadByte"

const WriteBoolFunction = "WriteBool"
const WriteIntFunction = "WriteUint64"
const WriteStringFunction = "WriteString"
const WriteByteFunction = "WriteByte"

type PrimitiveFunctions struct {
	ReadFunction  string
	WriteFunction string
	WriteCast     string
	ReadCast      string
}

// all integer types are read as uint64 then casted to the appropriate type
var Uint64PrimitiveFunctions = PrimitiveFunctions{
	ReadFunction:  ReadIntFunction,
	WriteFunction: WriteIntFunction,
}
var IntPrimitiveFunctions = PrimitiveFunctions{
	ReadFunction:  ReadIntFunction,
	WriteFunction: WriteIntFunction,
	WriteCast:     "uint64",
}

/*
The official implementation treats uint8s very strangely.  Slices of them are treated like []bytes (and rightly so because byte is just an aliased uint8), but individual uint8s are treated like other integer types (Uint64s).
We replicate this quirk.
*/
var UInt8SlicePrimitiveFunctions = PrimitiveFunctions{
	ReadFunction:  ReadByteFunction,
	WriteFunction: WriteByteFunction,
	WriteCast:     "uint8",
}

var BoolPrimitiveFunction = PrimitiveFunctions{ReadFunction: ReadBoolFunction, WriteFunction: WriteBoolFunction}
var StringPrimitiveFunction = PrimitiveFunctions{ReadFunction: ReadStringFunction, WriteFunction: WriteStringFunction}
var BytePrimitiveFunction = PrimitiveFunctions{ReadFunction: ReadByteFunction, WriteFunction: WriteByteFunction}

var supportedPrimitives = map[string]PrimitiveFunctions{
	"bool":   BoolPrimitiveFunction,
	"string": StringPrimitiveFunction,
	"byte":   BytePrimitiveFunction,
	"int":    IntPrimitiveFunctions,
	"int8":   IntPrimitiveFunctions,
	"int16":  IntPrimitiveFunctions,
	"int32":  IntPrimitiveFunctions,
	"int64":  IntPrimitiveFunctions,
	"uint":   IntPrimitiveFunctions,
	"uint8":  IntPrimitiveFunctions,
	"uint16": IntPrimitiveFunctions,
	"uint32": IntPrimitiveFunctions,
	"uint64": Uint64PrimitiveFunctions,
}

func isPrimitiveString(t string) bool {
	_, ok := supportedPrimitives[t]
	return ok
}
