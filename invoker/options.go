package invoker

type Options struct {
	Name          string
	InvokeOptions InvokeOptions
}

type Option func(*Options)

type InvokeOptions struct {
	Interceptors []Interceptor
}

type InvokeOption func(*InvokeOptions)

func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

func Intercept(i ...Interceptor) Option {
	return func(o *Options) {
		o.InvokeOptions.Interceptors = append(o.InvokeOptions.Interceptors, i...)
	}
}

func WithInterceptor(i ...Interceptor) InvokeOption {
	return func(o *InvokeOptions) {
		o.Interceptors = append(o.Interceptors, i...)
	}
}
