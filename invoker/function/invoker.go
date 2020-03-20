package function

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/haormj/util/invoker"
	"github.com/haormj/util/log"
)

type Invoker struct {
	opts invoker.Options
	fn   interface{}
	name string
	t    reflect.Type
	v    reflect.Value
	ft   funcType
}

type funcType struct {
	funcName   string
	Func       reflect.Type
	CtxType    reflect.Type
	InputType  reflect.Type
	OutputType reflect.Type
	ErrType    reflect.Type
}

func (f funcType) FuncName() string {
	return f.funcName
}

func (f funcType) In() []reflect.Type {
	return []reflect.Type{f.CtxType, f.InputType, f.OutputType}
}

func (f funcType) Out() []reflect.Type {
	return []reflect.Type{f.ErrType}
}

// Precompute the reflect type for error. Can't use error directly
// because Typeof takes an empty interface value. This is annoying.
var (
	typeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()
	typeOfError   = reflect.TypeOf((*error)(nil)).Elem()
)

func suitableFunc(typ reflect.Type, name string) (funcType, error) {
	mtype := typ
	mname := name
	var ft funcType
	// Method needs four ins: receiver, ctx, input, output
	if mtype.NumIn() != 3 {
		s := "method " + mname + " has wrong number of ins:" + string(mtype.NumIn())
		log.Error(s)
		return ft, errors.New(s)
	}
	// First arg must be context.Context
	ctxType := mtype.In(0)
	if ctxType != typeOfContext {
		s := "first param must be context.Context"
		log.Error(s)
		return ft, errors.New(s)
	}
	// Second arg need not be a pointer.
	inputType := mtype.In(1)
	if !isExportedOrBuiltinType(inputType) {
		s := mname + "argument type not exported:" + inputType.String()
		log.Error(s)
		return ft, errors.New(s)
	}
	// Output type must be exported.
	outputType := mtype.In(2)
	if outputType.Kind() != reflect.Ptr {
		s := "method " + mname + " reply type not a pointer:" + outputType.String()
		log.Error(s)
		return ft, errors.New(s)
	}
	// Output type must be exported.
	if !isExportedOrBuiltinType(outputType) {
		s := "method " + mname + " reply type not exported:" + outputType.String()
		log.Error(s)
		return ft, errors.New(s)
	}
	// Method needs one out.
	if mtype.NumOut() != 1 {
		s := "method " + mname + " has wrong number of outs:" + string(mtype.NumOut())
		log.Error(s)
		return ft, errors.New(s)
	}
	// The return type of the method must be error.
	errType := mtype.Out(0)
	if errType != typeOfError {
		s := "method" + mname + "returns" + errType.String() + "not error"
		log.Error(s)
		return ft, errors.New(s)
	}
	ft = funcType{
		funcName:   name,
		Func:       mtype,
		CtxType:    ctxType,
		InputType:  inputType,
		OutputType: outputType,
		ErrType:    errType,
	}
	return ft, nil
}

// Is this an exported - upper case - name?
func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}

func NewInvoker(fn interface{}, opts ...invoker.Option) invoker.Invoker {
	options := invoker.Options{}
	for _, o := range opts {
		o(&options)
	}

	inv := &Invoker{
		opts: options,
		fn:   fn,
	}

	return inv
}

func (d *Invoker) Name() string {
	return d.name
}

func (d *Invoker) Init(opts ...invoker.Option) error {
	for _, o := range opts {
		o(&d.opts)
	}
	t := reflect.TypeOf(d.fn)
	if t.Kind() != reflect.Func {
		s := t.Kind().String() + " is not func"
		log.Error(s)
		return errors.New(s)
	}
	d.t = t
	d.v = reflect.ValueOf(d.fn)
	funcName := runtime.FuncForPC(d.v.Pointer()).Name()
	if strings.Contains(funcName, ".") {
		splits := strings.Split(funcName, ".")
		funcName = splits[len(splits)-1]
	}

	if len(d.opts.Name) != 0 {
		d.name = d.opts.Name
	} else {
		d.name = funcName
	}

	ft, err := suitableFunc(t, funcName)
	if err != nil {
		return err
	}
	d.ft = ft
	return nil
}

func (f *Invoker) invoke(c context.Context, mi invoker.Message,
	opts ...invoker.InvokeOption) (mo invoker.Message, err error) {
	mo = invoker.NewMessage()
	if f.ft.FuncName() != mi.FuncName() {
		s := mi.FuncName() + " not find"
		log.Error(s)
		return mo, errors.New(s)
	}
	params := mi.Parameters()
	if len(params) != 3 {
		s := fmt.Sprintf("parameters must be 3, now %d", len(params))
		log.Error(s)
		return mo, errors.New(s)
	}
	if !reflect.TypeOf(params[0]).Implements(typeOfContext) {
		s := "first parameter must be context"
		log.Error(s)
		return mo, errors.New(s)
	}
	ctx := params[0].(context.Context)
	for k, v := range mi.Attachments() {
		ctx = context.WithValue(ctx, k, v)
	}
	params[0] = ctx
	in := []reflect.Value{}
	for _, p := range params {
		in = append(in, reflect.ValueOf(p))
	}
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			log.Error(err)
			return
		}
	}()
	out := f.v.Call(in)
	p := make([]interface{}, 0)
	for _, o := range out {
		p = append(p, o.Interface())
	}
	mo.SetParameters(p)
	return mo, nil
}

func (f *Invoker) Invoke(c context.Context, mi invoker.Message,
	opts ...invoker.InvokeOption) (invoker.Message, error) {
	invokeOptions := f.opts.InvokeOptions
	for _, o := range opts {
		o(&invokeOptions)
	}

	inv := f.invoke
	for i := len(invokeOptions.Interceptors); i > 0; i-- {
		inv = invokeOptions.Interceptors[i-1](inv)
	}

	return inv(c, mi, opts...)
}

func (f *Invoker) Function(string) (invoker.Function, error) {
	return f.ft, nil
}

func (d *Invoker) Functions() []invoker.Function {
	return []invoker.Function{d.ft}
}

func (d *Invoker) String() string {
	return "func"
}
