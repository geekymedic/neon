package validator

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
)

type OpFn func(*Field) error

// TODO: add more type check
func checkTypeValidate(arg interface{}, tag string) (reflect.Value, error) {
	var val = reflect.ValueOf(arg)
	switch val.Kind() {
	case reflect.Func, reflect.Chan, reflect.Interface, reflect.Invalid,
		reflect.Complex64, reflect.Complex128, reflect.Uintptr, reflect.UnsafePointer:
		return val, fmt.Errorf("tag: %v, not support %v's type", tag, val.Kind())
	default:
	}

	return val, nil
}

func isType(x interface{}, kind reflect.Kind, tag string) error {
	var val = reflect.ValueOf(x)
	if val.Kind() == kind {
		return nil
	}
	return fmt.Errorf("tag: %v, expect type: %v, actual: %v", tag, kind, val.Kind())
}

func eqType(ftype reflect.Kind, arg reflect.Kind) bool { return ftype == arg }

func neqKind(ftype reflect.Kind, arg reflect.Kind, tag string) error {
	return fmt.Errorf("tag: %v, expect type: %v, actual type: %v, not equal", tag, ftype, arg)
}

// ParseValidatorField parse ast file, then returns a slice
func ParseValidatorAst(structName, filename string) ([]string, error) {
	var (
		fields []string
		err    error
	)
	var fn leafFn = func(ctx context.Context, linkName string, _ ast.Expr) {
		fields = append(fields, linkName)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	for _, decl := range f.Decls {
		genSpec, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, typeSpec := range genSpec.Specs {
			if spec, ok := typeSpec.(*ast.TypeSpec); ok {
				if spec.Name.Name == structName {
					walkTypeSpec(context.TODO(), "", typeSpec, fn)
				}
			}
		}
	}

	return fields, err
}

type Field struct {
	Tag   string
	Value interface{}
}

func NewValidator() *Validator {
	return &Validator{err: nil}
}

type Validator struct {
	err error
}

func (m *Validator) Validate(field *Field, opFn OpFn) *Validator {
	m.err = opFn(field)
	return m
}

func (m *Validator) Ok() bool {
	return m.err == nil
}

func (m *Validator) Err() error { return m.err }

/////////////////////////////////////////////////
// ast parser
type leafFn func(ctx context.Context, linkName string, leafIdent ast.Expr)

func doLeafFns(ctx context.Context, linkName string, leafIdent ast.Expr, leafFns ...leafFn) {
	for _, leafFn := range leafFns {
		leafFn(ctx, linkName, leafIdent)
	}
}

func walkTypeSpec(ctx context.Context, parentName string, typeSpec ast.Spec, leafFns ...leafFn) {
	if typeSpec == nil {
		return
	}
	// get type
	spec, ok := typeSpec.(*ast.TypeSpec)
	if !ok {
		return
	}
	identName := spec.Name.Name // type IdentName srtruct, skip it
	structType, ok := spec.Type.(*ast.StructType)
	if !ok {
		return
	}
	if structType.Fields == nil {
		return
	}
	if parentName == "" {
		parentName = identName
	}
	if leafFns == nil {
		leafFns = []leafFn{}
	}
	for _, field := range structType.Fields.List {
		walkAstField(ctx, parentName, field, leafFns...)
	}
}

func walkAstField(ctx context.Context, parentName string, astField *ast.Field, leafFns ...leafFn) {
	if leafFns == nil {
		leafFns = []leafFn{}
	}
	if astField.Names == nil {
		// anonymous
		// eg: type A struct {Location} or type A struct {*Location}
		// A.Location == astField.type
		parentName = fmt.Sprintf("%s.%s", parentName, anonymousFieldName(astField.Type))
	} else {
		// variable
		// eg: type A struct {L Location}
		// A.L == astField.Names[0].Name
		parentName = fmt.Sprintf("%s.%s", parentName, astField.Names[0].Name)
	}

	walkExpr(ctx, parentName, astField.Type, leafFns...)
}

func anonymousFieldName(expr ast.Expr) string {
	switch realType := expr.(type) {
	case *ast.Ident:
		return realType.Name
	case *ast.StarExpr:
		return realType.X.(*ast.Ident).Name
	default:
		panic("unimplement")
	}
}

func walkExpr(ctx context.Context, parentName string, expr ast.Expr, leafFns ...leafFn) {
	if leafFns == nil {
		leafFns = []leafFn{}
	}
	switch realType := expr.(type) {
	case *ast.Ident:
		if obj := realType.Obj; obj != nil && obj.Decl != nil {
			spec, ok := obj.Decl.(*ast.TypeSpec)
			if ok {
				walkTypeSpec(ctx, parentName, spec, leafFns...)
			}
		} else {
			doLeafFns(ctx, parentName, realType, leafFns...)
		}
	case *ast.MapType: // terminal it, skip next level
		doLeafFns(ctx, parentName, realType, leafFns...)
		//walkExpr(ctx, parentName, realType.Value, leafFns...)
	case *ast.ArrayType: // terminal it, skip next level
		doLeafFns(ctx, parentName, realType, leafFns...)
		//walkExpr(ctx, parentName, realType.Elt, leafFns...)
	case *ast.StarExpr: // pointer type
		walkExpr(ctx, parentName, realType.X, leafFns...)
	case *ast.StructType:
		for _, field := range realType.Fields.List {
			walkAstField(ctx, parentName, field, leafFns...)
		}
	default:
		panic("unimplement")
	}
}
