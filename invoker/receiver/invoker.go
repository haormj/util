package receiver

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf8"

	"github.com/haormj/util/invoker"
	"github.com/haormj/util/log"
)

type ReceiverInvoker struct {
	receiver interface{}
	name     string                 // name of service
	rcvr     reflect.Value          // receiver of methods for the service
	typ      reflect.Type           // type of the receiver
	method   map[string]*methodType // registered methods

	opts invoker.Options
}

type methodType struct {
	method     reflect.Method
	CtxType    reflect.Type
	InputType  reflect.Type
	OutputType reflect.Type
	ErrType    reflect.Type
}

func (m *methodType) FuncName() string {
	return m.method.Name
}

func (m *methodType) In() []reflect.Type {
	return []reflect.Type{m.CtxType, m.InputType, m.OutputType}
}

func (m *methodType) Out() []reflect.Type {
	return []reflect.Type{m.ErrType}
}

// Precompute the reflect type for error. Can't use error directly
// because Typeof takes an empty interface value. This is annoying.
var (
	typeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()
	typeOfError   = reflect.TypeOf((*error)(nil)).Elem()
)

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

func (r *ReceiverInvoker) register(rcvr interface{}, name string, useName bool) error {
	r.typ = reflect.TypeOf(rcvr)
	r.rcvr = reflect.ValueOf(rcvr)
	sname := reflect.Indirect(r.rcvr).Type().Name()
	if useName {
		sname = name
	}
	if sname == "" {
		s := "rpc.Register: no service name for type " + r.typ.String()
		log.Error(s)
		return errors.New(s)
	}
	if !isExported(sname) && !useName {
		s := "rpc.Register: type " + sname + " is not exported"
		log.Error(s)
		return errors.New(s)
	}
	r.name = sname

	// Install the methods
	r.method = suitableMethods(r.typ, true)

	if len(r.method) == 0 {
		str := ""

		// To help the user, see if a pointer receiver would work.
		method := suitableMethods(reflect.PtrTo(r.typ), false)
		if len(method) != 0 {
			str = "rpc.Register: type " + sname + " has no exported methods of suitable type (hint: pass a pointer to value of that type)"
		} else {
			str = "rpc.Register: type " + sname + " has no exported methods of suitable type"
		}
		log.Error(str)
		return errors.New(str)
	}
	return nil
}

// suitableMethods returns suitable Rpc methods of typ, it will report
// error using log if reportErr is true.
func suitableMethods(typ reflect.Type, reportErr bool) map[string]*methodType {
	methods := make(map[string]*methodType)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name
		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}
		// Method needs four ins: receiver, ctx, input, output
		if mtype.NumIn() != 4 {
			if reportErr {
				log.Warn("method ", mname, " has wrong number of ins:", mtype.NumIn())
			}
			continue
		}
		// First arg must be context.Context
		ctxType := mtype.In(1)
		if ctxType != typeOfContext {
			continue
		}
		// Second arg need not be a pointer.
		inputType := mtype.In(2)
		if !isExportedOrBuiltinType(inputType) {
			if reportErr {
				log.Warn(mname, "argument type not exported:", inputType)
			}
			continue
		}
		// Output type must be exported.
		outputType := mtype.In(3)
		if outputType.Kind() != reflect.Ptr {
			if reportErr {
				log.Warn("method ", mname, " reply type not a pointer:", outputType)
			}
			continue
		}
		// Output type must be exported.
		if !isExportedOrBuiltinType(outputType) {
			if reportErr {
				log.Warn("method ", mname, " reply type not exported:", outputType)
			}
			continue
		}
		// Method needs one out.
		if mtype.NumOut() != 1 {
			if reportErr {
				log.Warn("method ", mname, " has wrong number of outs:", mtype.NumOut())
			}
			continue
		}
		// The return type of the method must be error.
		errType := mtype.Out(0)
		if errType != typeOfError {
			if reportErr {
				log.Warn("method", mname, "returns", errType.String(), "not error")
			}
			continue
		}
		methods[mname] = &methodType{
			method:     method,
			CtxType:    ctxType,
			InputType:  inputType,
			OutputType: outputType,
			ErrType:    errType,
		}
	}
	return methods
}

func NewInvoker(receiver interface{}, opts ...invoker.Option) invoker.Invoker {
	options := invoker.Options{}
	for _, o := range opts {
		o(&options)
	}

	inv := &ReceiverInvoker{
		receiver: receiver,
		opts:     options,
	}

	return inv
}

func (r *ReceiverInvoker) Init(opts ...invoker.Option) error {
	for _, o := range opts {
		o(&r.opts)
	}

	if r.receiver == nil {
		return errors.New("receiver is nil")
	}

	useName := false
	name := ""
	if len(r.opts.Name) != 0 {
		name = r.opts.Name
		useName = true
	}
	if err := r.register(r.receiver, name, useName); err != nil {
		return err
	}
	return nil
}

func (r *ReceiverInvoker) invoke(c context.Context, mi invoker.Message,
	opts ...invoker.InvokeOption) (mo invoker.Message, err error) {
	mo = invoker.NewMessage()
	mtype, ok := r.method[mi.FuncName()]
	if !ok {
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
	in := []reflect.Value{r.rcvr}
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
	out := mtype.method.Func.Call(in)
	p := make([]interface{}, 0)
	for _, o := range out {
		p = append(p, o.Interface())
	}
	mo.SetParameters(p)
	return mo, nil
}

func (r *ReceiverInvoker) Name() string {
	return r.name
}

func (r *ReceiverInvoker) Invoke(c context.Context, m invoker.Message,
	opts ...invoker.InvokeOption) (invoker.Message, error) {
	invokeOptions := r.opts.InvokeOptions
	for _, o := range opts {
		o(&invokeOptions)
	}

	inv := r.invoke
	for i := len(invokeOptions.Interceptors); i > 0; i-- {
		inv = invokeOptions.Interceptors[i-1](inv)
	}

	return inv(c, m, opts...)
}

func (r *ReceiverInvoker) Function(n string) (invoker.Function, error) {
	mtype, ok := r.method[n]
	if !ok {
		s := n + " not find"
		log.Error(s)
		return nil, errors.New(s)
	}
	return mtype, nil
}

func (r *ReceiverInvoker) Functions() []invoker.Function {
	fs := make([]invoker.Function, 0)
	for _, f := range r.method {
		fs = append(fs, f)
	}
	return fs
}

func (r *ReceiverInvoker) String() string {
	return "rcvr"
}
