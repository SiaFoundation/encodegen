package codegen

const ReadBoolFunction = "ReadBool"
const ReadIntFunction = "ReadUint64"
const ReadStringFunction = "ReadPrefixedBytes"
const ReadByteFunction = "ReadByte"

const WriteBoolFunction = "WriteBool"
const WriteIntFunction = "WriteUint64"
const WriteStringFunction = "WritePrefixedBytes"
const WriteByteFunction = "WriteByte"

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

var BoolPrimitiveFunction = PrimitiveFunctions{ReadFunction: ReadBoolFunction, WriteFunction: WriteBoolFunction, WriteCast: "bool"}
var StringPrimitiveFunction = PrimitiveFunctions{ReadFunction: ReadStringFunction, WriteFunction: WriteStringFunction, WriteCast: "[]byte"}
var BytePrimitiveFunction = PrimitiveFunctions{ReadFunction: ReadByteFunction, WriteFunction: WriteByteFunction, WriteCast: "byte"}

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
	"uint64": IntPrimitiveFunctions,

	"*bool":   BoolPrimitiveFunction,
	"*string": StringPrimitiveFunction,
	"*byte":   BytePrimitiveFunction,
	"*int":    IntPrimitiveFunctions,
	"*int8":   IntPrimitiveFunctions,
	"*int16":  IntPrimitiveFunctions,
	"*int32":  IntPrimitiveFunctions,
	"*int64":  IntPrimitiveFunctions,
	"*uint":   IntPrimitiveFunctions,
	"*uint8":  IntPrimitiveFunctions,
	"*uint16": IntPrimitiveFunctions,
	"*uint32": IntPrimitiveFunctions,
	"*uint64": IntPrimitiveFunctions,
}

var supportedPrimitivesArray = map[string]PrimitiveFunctions{
	"[]bool":   BoolPrimitiveFunction,
	"[]string": StringPrimitiveFunction,
	"[]byte":   BytePrimitiveFunction,
	"[]int":    IntPrimitiveFunctions,
	"[]int8":   IntPrimitiveFunctions,
	"[]int16":  IntPrimitiveFunctions,
	"[]int32":  IntPrimitiveFunctions,
	"[]int64":  IntPrimitiveFunctions,
	"[]uint":   IntPrimitiveFunctions,
	"[]uint8":  IntPrimitiveFunctions,
	"[]uint16": IntPrimitiveFunctions,
	"[]uint32": IntPrimitiveFunctions,
	"[]uint64": IntPrimitiveFunctions,

	"*[]bool":   BoolPrimitiveFunction,
	"*[]string": StringPrimitiveFunction,
	"*[]byte":   BytePrimitiveFunction,
	"*[]int":    IntPrimitiveFunctions,
	"*[]int8":   IntPrimitiveFunctions,
	"*[]int16":  IntPrimitiveFunctions,
	"*[]int32":  IntPrimitiveFunctions,
	"*[]int64":  IntPrimitiveFunctions,
	"*[]uint":   IntPrimitiveFunctions,
	"*[]uint8":  IntPrimitiveFunctions,
	"*[]uint16": IntPrimitiveFunctions,
	"*[]uint32": IntPrimitiveFunctions,
	"*[]uint64": IntPrimitiveFunctions,

	"[]*bool":   BoolPrimitiveFunction,
	"[]*string": StringPrimitiveFunction,
	"[]*byte":   BytePrimitiveFunction,
	"[]*int":    IntPrimitiveFunctions,
	"[]*int8":   IntPrimitiveFunctions,
	"[]*int16":  IntPrimitiveFunctions,
	"[]*int32":  IntPrimitiveFunctions,
	"[]*int64":  IntPrimitiveFunctions,
	"[]*uint":   IntPrimitiveFunctions,
	"[]*uint8":  IntPrimitiveFunctions,
	"[]*uint16": IntPrimitiveFunctions,
	"[]*uint32": IntPrimitiveFunctions,
	"[]*uint64": IntPrimitiveFunctions,

	"*[]*bool":   BoolPrimitiveFunction,
	"*[]*string": StringPrimitiveFunction,
	"*[]*byte":   BytePrimitiveFunction,
	"*[]*int":    IntPrimitiveFunctions,
	"*[]*int8":   IntPrimitiveFunctions,
	"*[]*int16":  IntPrimitiveFunctions,
	"*[]*int32":  IntPrimitiveFunctions,
	"*[]*int64":  IntPrimitiveFunctions,
	"*[]*uint":   IntPrimitiveFunctions,
	"*[]*uint8":  IntPrimitiveFunctions,
	"*[]*uint16": IntPrimitiveFunctions,
	"*[]*uint32": IntPrimitiveFunctions,
	"*[]*uint64": IntPrimitiveFunctions,
}

func isPrimitiveString(t string) bool {
	_, ok := supportedPrimitives[t]
	return ok
}

func isPrimitiveArrayString(t string) bool {
	_, ok := supportedPrimitivesArray[t]
	return ok
}
