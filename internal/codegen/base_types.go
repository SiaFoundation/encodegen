package codegen

const ReadBoolFunction = "readBool"
const ReadIntFunction = "readUint64"
const ReadStringFunction = "readPrefixedBytes"
const ReadByteFunction = "readByte"

const WriteBoolFunction = "writeBool"
const WriteIntFunction = "writeUint64"
const WriteStringFunction = "writePrefixedBytes"
const WriteByteFunction = "writeByte"

type PrimitiveFunctions struct {
	ReadFunction  string
	WriteFunction string
	WriteCast     string
}

var IntPrimitiveFunctions = PrimitiveFunctions{
	ReadFunction:  ReadIntFunction,
	WriteFunction: WriteIntFunction,
	WriteCast:     "uint64",
}

var supportedPrimitives = map[string]PrimitiveFunctions{
	"bool":   {ReadFunction: ReadBoolFunction, WriteFunction: WriteBoolFunction, WriteCast: "bool"},
	"string": {ReadFunction: ReadStringFunction, WriteFunction: WriteStringFunction, WriteCast: "[]byte"},
	"byte":   {ReadFunction: ReadByteFunction, WriteFunction: WriteByteFunction, WriteCast: "byte"},
	"int":    IntPrimitiveFunctions,
	"int8":   IntPrimitiveFunctions,
	"int16":  IntPrimitiveFunctions,
	"int32":  IntPrimitiveFunctions,
	"int64":  IntPrimitiveFunctions,
	"uint":   IntPrimitiveFunctions,
	"uint8":  IntPrimitiveFunctions,
	"uint16": IntPrimitiveFunctions,
	"uint32": IntPrimitiveFunctions,
	"uint64": IntPrimitiveFunctions,
}

func isPrimitiveString(t string) bool {
	_, ok := supportedPrimitives[t]
	return ok
}
