package invoker

import (
	"context"
	"reflect"
	"sync"
)

type Invoker interface {
	Name() string
	Init(...Option) error
	Invoke(context.Context, Message, ...InvokeOption) (Message, error)
	Function(string) (Function, error)
	Functions() []Function
	String() string
}

type Message interface {
	FuncName() string
	SetFuncName(string)
	Parameters() []interface{}
	SetParameters([]interface{})
	Attachments() map[string]string
	Attachment(string) (string, bool)
	SetAttachment(string, string)
}

type Function interface {
	FuncName() string
	In() []reflect.Type
	Out() []reflect.Type
}

type InvokeFunc func(context.Context, Message, ...InvokeOption) (Message, error)

type Interceptor func(InvokeFunc) InvokeFunc

type message struct {
	sync.RWMutex
	funcName    string
	parameters  []interface{}
	attachments map[string]string
}

func NewMessage() Message {
	r := &message{
		parameters:  make([]interface{}, 0),
		attachments: make(map[string]string),
	}
	return r
}

func (r *message) FuncName() string {
	r.RLock()
	f := r.funcName
	r.RUnlock()
	return f
}

func (r *message) SetFuncName(f string) {
	r.Lock()
	r.funcName = f
	r.Unlock()
}

func (r *message) Parameters() []interface{} {
	r.RLock()
	p := r.parameters
	r.RUnlock()
	return p
}

func (r *message) SetParameters(p []interface{}) {
	r.Lock()
	r.parameters = p
	r.Unlock()
}

func (r *message) Attachments() map[string]string {
	m := make(map[string]string)
	r.RLock()
	for k, v := range r.attachments {
		m[k] = v
	}
	r.RUnlock()
	return m
}
func (r *message) Attachment(k string) (string, bool) {
	r.RLock()
	v, ok := r.attachments[k]
	r.RUnlock()
	return v, ok
}
func (r *message) SetAttachment(k string, v string) {
	r.Lock()
	r.attachments[k] = v
	r.Unlock()
}
