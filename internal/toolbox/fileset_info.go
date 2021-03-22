package toolbox

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

//FieldInfo represents a filed info
type FieldInfo struct {
	Name                 string
	TypeName             string
	ComponentType        string
	IsPointerComponent   bool
	KeyTypeName          string
	ValueTypeName        string
	TypePackage          string
	IsMap                bool
	IsChannel            bool
	IsSlice              bool
	IsPointer            bool
	Tag                  string
	Comment              string
	IsVariant            bool
	IsFixed              bool
	FixedSize            int
	AnonymousChildFields []*FieldInfo
}

//NewFunctionInfoFromField creates a new function info.
func NewFunctionInfoFromField(field *ast.Field, owner *FileInfo) *FunctionInfo {
	result := &FunctionInfo{
		Name:            "",
		ParameterFields: make([]*FieldInfo, 0),
		ResultsFields:   make([]*FieldInfo, 0),
	}
	if len(field.Names) > 0 {
		result.Name = field.Names[0].Name
	}

	if funcType, ok := field.Type.(*ast.FuncType); ok {
		if funcType.Params != nil && len(funcType.Params.List) > 0 {
			result.ParameterFields = fieldListToFieldInfoSlice(funcType.Params)
		}
		if funcType.Results != nil && len(funcType.Results.List) > 0 {
			result.ResultsFields = fieldListToFieldInfoSlice(funcType.Results)
		}
		var names = make(map[string]bool)
		for _, param := range result.ParameterFields {
			if strings.Contains(strings.ToLower(param.TypeName), strings.ToLower(param.Name)) {
				name := matchLastNameSegment(param.TypeName)
				if _, has := names[name]; has {
					continue
				}
				names[name] = true
				param.Name = name
			}
		}
	}
	return result
}

func matchLastNameSegment(name string) string {
	var result = make([]byte, 0)
	for i := len(name) - 1; i >= 0; i-- {
		aChar := string(name[i : i+1])
		if aChar != "." {
			result = append(result, byte(aChar[0]))
		}
		if strings.ToUpper(aChar) == aChar || aChar == "." {
			ReverseSlice(result)
			return string(result)
		}
	}
	return name
}

//NewFieldInfo creates a new field info.
func NewFieldInfo(field *ast.Field) *FieldInfo {
	return NewFieldInfoByIndex(field, 0)
}

//NewFieldInfoByIndex creates a new field info.
func NewFieldInfoByIndex(field *ast.Field, index int) *FieldInfo {
	result := &FieldInfo{
		Name:     "",
		TypeName: types.ExprString(field.Type),
	}

	if len(field.Names) > 0 {
		result.Name = field.Names[index].Name
	}
	_, result.IsMap = field.Type.(*ast.MapType)
	var arrayType *ast.ArrayType
	if arrayType, result.IsSlice = field.Type.(*ast.ArrayType); result.IsSlice {
		switch x := arrayType.Elt.(type) {
		case *ast.Ident:
			result.ComponentType = x.Name
		case *ast.StarExpr:
			switch y := x.X.(type) {
			case *ast.Ident:
				result.ComponentType = y.Name
			case *ast.SelectorExpr:
				result.ComponentType = y.X.(*ast.Ident).Name + "." + y.Sel.Name
			}
			result.IsPointerComponent = true
		case *ast.SelectorExpr:
			result.ComponentType = x.X.(*ast.Ident).Name + "." + x.Sel.Name
		}
	}
	_, result.IsPointer = field.Type.(*ast.StarExpr)
	_, result.IsChannel = field.Type.(*ast.ChanType)
	if selector, ok := field.Type.(*ast.SelectorExpr); ok {
		result.TypePackage = types.ExprString(selector.X)
	}
	if result.IsPointer {
		if pointerExpr, casted := field.Type.(*ast.StarExpr); casted {
			if identExpr, ok := pointerExpr.X.(*ast.Ident); ok {
				result.TypeName = identExpr.Name
			}
		}
	} else if identExpr, ok := field.Type.(*ast.Ident); ok {
		result.TypeName = identExpr.Name
	}

	if field.Tag != nil {
		result.Tag = field.Tag.Value
	}
	if mapType, ok := field.Type.(*ast.MapType); ok {
		result.KeyTypeName = types.ExprString(mapType.Key)
		result.ValueTypeName = types.ExprString(mapType.Value)
	}

	if strings.Contains(result.TypeName, "...") {
		result.IsVariant = true
		result.TypeName = strings.Replace(result.TypeName, "...", "[]", 1)
	}

	if index := strings.Index(result.TypeName, "."); index != -1 {
		from := 0
		if result.IsPointer {
			from = 1
		}
		result.TypePackage = string(result.TypeName[from:index])
	}
	// fmt.Printf("Result: {%+v}\n", result)

	return result
}

//FunctionInfo represents a function info
type FunctionInfo struct {
	Name             string
	ReceiverTypeName string
	ParameterFields  []*FieldInfo
	ResultsFields    []*FieldInfo
	*FileInfo
}

//NewFunctionInfo create a new function
func NewFunctionInfo(funcDeclaration *ast.FuncDecl, owner *FileInfo) *FunctionInfo {
	result := &FunctionInfo{
		Name:            "",
		ParameterFields: make([]*FieldInfo, 0),
		ResultsFields:   make([]*FieldInfo, 0),
	}

	if funcDeclaration.Name != nil {
		result.Name = funcDeclaration.Name.Name
	}
	if funcDeclaration.Recv != nil {
		receiverType := funcDeclaration.Recv.List[0].Type
		if ident, ok := receiverType.(*ast.Ident); ok {
			result.ReceiverTypeName = ident.Name
		} else if startExpr, ok := receiverType.(*ast.StarExpr); ok {
			if ident, ok := startExpr.X.(*ast.Ident); ok {
				result.ReceiverTypeName = ident.Name
			}
		}
	}
	return result
}

//TypeInfo represents a struct info
type TypeInfo struct {
	Name                   string
	Package                string
	FileName               string
	Comment                string
	IsSlice                bool
	IsStruct               bool
	IsInterface            bool
	IsDerived              bool
	ComponentType          string
	IsPointerComponentType bool
	Derived                string
	IsFixed                bool
	FixedSize              int
	Settings               map[string]string
	fields                 []*FieldInfo
	indexedField           map[string]*FieldInfo
	receivers              []*FunctionInfo
	indexedReceiver        map[string]*FunctionInfo
	rcv                    *FunctionInfo
	// coveredTypes map[string]bool
}

//AddFields appends fileds to structinfo
func (s *TypeInfo) AddFields(fields ...*FieldInfo) {
	if s == nil {
		return
	}
	for _, field := range fields {
		// if s.coveredTypes[field.Name] {
		// 	fmt.Printf("(%+v) Already have a %s, continuing\n", s, field.Name)
		// 	continue
		// }
		// fmt.Printf("(%+v) Added a %s\n", s, field.Name)
		// s.coveredTypes[field.Name] = true

		s.fields = append(s.fields, field)
		s.indexedField[field.Name] = field
	}
}

//Field returns filedinfo for supplied file name
func (s *TypeInfo) Field(name string) *FieldInfo {
	return s.indexedField[name]
}

//Fields returns all fields
func (s *TypeInfo) Fields() []*FieldInfo {
	return s.fields
}

//HasField returns true if struct has passed in field.
func (s *TypeInfo) HasField(name string) bool {
	_, found := s.indexedField[name]
	return found
}

//Receivers returns struct functions
func (s *TypeInfo) Receivers() []*FunctionInfo {
	return s.receivers
}

//Receiver returns receiver for passed in name
func (s *TypeInfo) Receiver(name string) *FunctionInfo {
	return s.indexedReceiver[name]
}

//HasReceiver returns true if receiver is defined for struct
func (s *TypeInfo) HasReceiver(name string) bool {
	_, found := s.indexedReceiver[name]
	return found
}

//AddReceivers adds receiver for the struct
func (s *TypeInfo) AddReceivers(receivers ...*FunctionInfo) {
	s.receivers = append(s.receivers, receivers...)
	for _, receiver := range receivers {
		s.indexedReceiver[receiver.Name] = receiver
	}
}

//NewTypeInfo creates a new struct info
func NewTypeInfo(name string) *TypeInfo {
	return &TypeInfo{
		Name:            name,
		fields:          make([]*FieldInfo, 0),
		receivers:       make([]*FunctionInfo, 0),
		indexedReceiver: make(map[string]*FunctionInfo),
		indexedField:    make(map[string]*FieldInfo),
		Settings:        make(map[string]string),
		// coveredTypes: make(map[string]bool),
	}
}

//FileInfo represent hold definition about all defined types and its receivers in a file
type FileInfo struct {
	basePath            string
	filename            string
	types               map[string]*TypeInfo
	functions           map[string][]*FunctionInfo
	packageName         string
	currentTypInfo      *TypeInfo
	fileSet             *token.FileSet
	currentFunctionInfo *FunctionInfo
	Imports             map[string]string
	coveredFields       map[*ast.Field]bool
}

//Type returns a type info for passed in name
func (f *FileInfo) Type(name string) *TypeInfo {
	return f.types[name]
}

//Type returns a struct info for passed in name
func (f *FileInfo) addFunction(funcion *FunctionInfo) {
	functions, found := f.functions[funcion.ReceiverTypeName]
	if !found {
		functions = make([]*FunctionInfo, 0)
		f.functions[funcion.ReceiverTypeName] = functions
	}
	f.functions[funcion.ReceiverTypeName] = append(f.functions[funcion.ReceiverTypeName], funcion)
}

//Types returns all struct info
func (f *FileInfo) Types() []*TypeInfo {
	var result = make([]*TypeInfo, 0)
	for _, v := range f.types {
		result = append(result, v)
	}
	return result
}

//HasType returns truc if struct info is defined in a file
func (f *FileInfo) HasType(name string) bool {
	_, found := f.types[name]
	return found
}

//toFieldInfoSlice converts filedList to FiledInfo slice.
func fieldListToFieldInfoSlice(source *ast.FieldList) []*FieldInfo {
	var result = make([]*FieldInfo, 0)
	if source == nil || len(source.List) == 0 {
		return result
	}
	for _, field := range source.List {
		for i := range field.Names {
			result = append(result, NewFieldInfoByIndex(field, i))
		}
	}
	return result
}

//toFunctionInfos convers filedList to function info slice.
func toFunctionInfos(source *ast.FieldList, owner *FileInfo) []*FunctionInfo {
	var result = make([]*FunctionInfo, 0)
	if source == nil || len(source.List) == 0 {
		return result
	}
	for _, field := range source.List {
		result = append(result, NewFunctionInfoFromField(field, owner))
	}
	return result
}

func (f *FileInfo) getFieldWithAnonymousChildren(field *ast.Field) *FieldInfo {
	var newField *FieldInfo
	switch value := field.Type.(type) {
	default:
		newField = NewFieldInfo(field)
	case *ast.StructType:
		newField = f.anonymousAddStructFields(field, value)
	case *ast.StarExpr:
		switch xTypeValue := value.X.(type) {
		default:
			newField = NewFieldInfo(field)
		case *ast.StructType:
			newField = f.anonymousAddStructFields(field, xTypeValue)
		}
	case *ast.ArrayType:
		switch eltTypeValue := value.Elt.(type) {
		default:
			newField = NewFieldInfo(field)
		case *ast.StarExpr:
			switch xTypeValue := eltTypeValue.X.(type) {
			default:
				newField = NewFieldInfo(field)
			case *ast.StructType:
				newField = f.anonymousAddStructFields(field, xTypeValue)
			}
		case *ast.StructType:
			newField = f.anonymousAddStructFields(field, eltTypeValue)
		}
		if value.Len != nil {
			lenType, ok := value.Len.(*ast.BasicLit)
			if !ok {
				return newField
			}

			lenValue, err := strconv.Atoi(lenType.Value)
			if err != nil {
				return newField
			}
			newField.IsFixed = true
			newField.FixedSize = lenValue
		}
	}
	return newField
}

func (f *FileInfo) anonymousAddStructFields(structField *ast.Field, structType *ast.StructType) *FieldInfo {
	var fields []*FieldInfo

	for _, field := range structType.Fields.List {
		f.coveredFields[field] = true
		fields = append(fields, f.getFieldWithAnonymousChildren(field))
	}

	newField := NewFieldInfo(structField)
	newField.AnonymousChildFields = fields
	return newField
}

//Visit visits ast node to extract struct details from the passed file
func (f *FileInfo) Visit(node ast.Node) ast.Visitor {
	if node != nil {
		// fmt.Printf("node %+v\n", node)
		// spew.Dump(node)

		switch value := node.(type) {
		case *ast.TypeSpec:
			typeName := value.Name.Name
			typeInfo := NewTypeInfo(typeName)
			typeInfo.Package = f.packageName
			typeInfo.FileName = f.filename
			switch typeValue := value.Type.(type) {
			case *ast.ArrayType:
				typeInfo.IsSlice = true
				if ident, ok := typeValue.Elt.(*ast.Ident); ok {
					typeInfo.ComponentType = ident.Name
				} else if startExpr, ok := typeValue.Elt.(*ast.StarExpr); ok {
					if ident, ok := startExpr.X.(*ast.Ident); ok {
						typeInfo.ComponentType = ident.Name
					}
					typeInfo.IsPointerComponentType = true
				}
				if typeValue.Len != nil {
					typeInfo.IsFixed = true
					lenType, ok := typeValue.Len.(*ast.BasicLit)
					if !ok {
						break
					}

					lenValue, err := strconv.Atoi(lenType.Value)
					if err != nil {
						break
					}
					typeInfo.FixedSize = lenValue
				}

			case *ast.StructType:
				typeInfo.IsStruct = true
			case *ast.InterfaceType:
				typeInfo.IsInterface = true
			case *ast.Ident:
				typeInfo.Derived = typeValue.Name
				typeInfo.IsDerived = true
			}
			f.currentTypInfo = typeInfo
			f.types[typeName] = typeInfo
		case *ast.Field:
			if len(value.Names) < 1 || f.coveredFields[value] {
				break
			}
			f.currentTypInfo.AddFields(f.getFieldWithAnonymousChildren(value))
			f.coveredFields[value] = true
		case *ast.FuncDecl:
			functionInfo := NewFunctionInfo(value, f)
			functionInfo.FileInfo = f
			f.currentFunctionInfo = functionInfo
			if len(functionInfo.ReceiverTypeName) > 0 {
				f.addFunction(functionInfo)
			}

		case *ast.FuncType:

			if f.currentFunctionInfo != nil {
				if value.Params != nil {
					f.currentFunctionInfo.ParameterFields = fieldListToFieldInfoSlice(value.Params)
				}

				if value.Results != nil {
					f.currentFunctionInfo.ResultsFields = fieldListToFieldInfoSlice(value.Results)
				}
				f.currentFunctionInfo = nil
			}
		case *ast.FieldList:
			if f.currentTypInfo != nil && f.currentTypInfo.IsInterface {
				f.currentTypInfo.receivers = toFunctionInfos(value, f)
				f.currentTypInfo = nil
			}
		case *ast.ImportSpec:
			if value.Name != nil && value.Name.String() != "" {
				f.Imports[value.Name.String()] = value.Path.Value
			} else {
				_, name := path.Split(value.Path.Value)
				name = strings.Replace(name, `"`, "", 2)
				f.Imports[name] = value.Path.Value
			}
		}
	}
	return f
}

//NewFileInfo creates a new file info.
func NewFileInfo(basePath, packageName, filename string, fileSet *token.FileSet) *FileInfo {
	result := &FileInfo{
		basePath:      basePath,
		filename:      filename,
		packageName:   packageName,
		types:         make(map[string]*TypeInfo),
		functions:     make(map[string][]*FunctionInfo),
		Imports:       make(map[string]string),
		fileSet:       fileSet,
		coveredFields: make(map[*ast.Field]bool),
	}
	return result
}

//FileSetInfo represents a fileset info storing information about go file with their struct definition
type FileSetInfo struct {
	files map[string]*FileInfo
}

//FileInfo returns fileinfo for supplied file name
func (f *FileSetInfo) FileInfo(name string) *FileInfo {
	return f.files[name]
}

//FilesInfo returns all files info.
func (f *FileSetInfo) FilesInfo() map[string]*FileInfo {
	return f.files
}

//Type returns type info for passed in type  name.
func (f *FileSetInfo) Type(name string) *TypeInfo {
	if pointerIndex := strings.LastIndex(name, "*"); pointerIndex != -1 {
		name = name[pointerIndex+1:]
	}
	for _, v := range f.files {
		if v.HasType(name) {
			return v.Type(name)
		}
	}
	return nil
}

//NewFileSetInfo creates a new fileset info
func NewFileSetInfo(baseDir string) (*FileSetInfo, error) {
	fileSet := token.NewFileSet()
	pkgs, err := parser.ParseDir(fileSet, baseDir, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse path %v: %v", baseDir, err)
	}

	var result = &FileSetInfo{
		files: make(map[string]*FileInfo),
	}
	for packageName, pkg := range pkgs {
		for filename, file := range pkg.Files {
			filename := filepath.Base(filename)
			fileInfo := NewFileInfo(baseDir, packageName, filename, fileSet)
			ast.Walk(fileInfo, file)
			result.files[filename] = fileInfo
		}
	}

	for _, fileInfo := range result.files {
		for k, functionsInfo := range fileInfo.functions {
			typeInfo := result.Type(k)
			if typeInfo != nil && typeInfo.IsStruct {
				typeInfo.AddReceivers(functionsInfo...)
			}
		}

	}
	return result, nil
}
