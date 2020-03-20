package dummy

import (
	"context"

	"github.com/haormj/util/invoker"
)

type Invoker struct {
	opts invoker.Options
}

func NewInvoker(opts ...invoker.Option) invoker.Invoker {
	options := invoker.Options{}
	for _, o := range opts {
		o(&options)
	}

	inv := &Invoker{
		opts: options,
	}

	return inv
}

func (d *Invoker) Name() string {
	return d.opts.Name
}

func (d *Invoker) Init(opts ...invoker.Option) error {
	for _, o := range opts {
		o(&d.opts)
	}
	return nil
}

func (f *Invoker) invoke(c context.Context, mi invoker.Message,
	opts ...invoker.InvokeOption) (invoker.Message, error) {
	return nil, nil
}

func (d *Invoker) Invoke(c context.Context, mi invoker.Message,
	opts ...invoker.InvokeOption) (invoker.Message, error) {
	invokeOptions := d.opts.InvokeOptions
	for _, o := range opts {
		o(&invokeOptions)
	}

	inv := d.invoke
	for i := len(invokeOptions.Interceptors); i > 0; i-- {
		inv = invokeOptions.Interceptors[i-1](inv)
	}

	return inv(c, mi, opts...)
}

func (d *Invoker) Function(string) (invoker.Function, error) {
	return nil, nil
}

func (d *Invoker) Functions() []invoker.Function {
	return nil
}

func (d *Invoker) String() string {
	return "dummy"
}
