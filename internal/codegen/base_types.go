package codegen

import (
	"unsafe"
)

const ReadBoolFunction = "ReadBool"
const ReadIntFunction = "ReadUint64"
const ReadStringFunction = "ReadString"
const ReadByteFunction = "ReadByte"

const WriteBoolFunction = "WriteBool"
const WriteIntFunction = "WriteUint64"
const WriteStringFunction = "WriteString"
const WriteByteFunction = "WriteByte"

const sizeofUint64  = int(unsafe.Sizeof(uint64(0)))
const sizeofByte  = int(unsafe.Sizeof(byte(0)))
const sizeofBool  = int(unsafe.Sizeof(bool(false)))

type PrimitiveFunctions struct {
	ReadFunction  string
	WriteFunction string
	WriteCast     string
	ReadCast      string
	ResetString   string
	ElementSize   int
}

// all integer types are read as uint64 then casted to the appropriate type
var Uint64PrimitiveFunctions = PrimitiveFunctions{
	ReadFunction:  ReadIntFunction,
	WriteFunction: WriteIntFunction,
	ResetString:   "0",
	ElementSize:   sizeofUint64,
}
var IntPrimitiveFunctions = PrimitiveFunctions{
	ReadFunction:  ReadIntFunction,
	WriteFunction: WriteIntFunction,
	WriteCast:     "uint64",
	ResetString:   "0",
	ElementSize:   sizeofUint64,
}

/*
The official implementation treats uint8s very strangely.  Slices of them are treated like []bytes (and rightly so because byte is just an aliased uint8), but individual uint8s are treated like other integer types (Uint64s).
We replicate this quirk.
*/
var UInt8SlicePrimitiveFunctions = PrimitiveFunctions{
	ReadFunction:  ReadByteFunction,
	WriteFunction: WriteByteFunction,
	WriteCast:     "uint8",
	ResetString:   "0",
	ElementSize:   sizeofByte,
}

var BoolPrimitiveFunction = PrimitiveFunctions{ReadFunction: ReadBoolFunction, WriteFunction: WriteBoolFunction, ResetString: "false", ElementSize: sizeofBool}
var StringPrimitiveFunction = PrimitiveFunctions{ReadFunction: ReadStringFunction, WriteFunction: WriteStringFunction, ResetString: `""`, ElementSize: sizeofByte}
var BytePrimitiveFunction = PrimitiveFunctions{ReadFunction: ReadByteFunction, WriteFunction: WriteByteFunction, ResetString: "0", ElementSize: sizeofByte}

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
