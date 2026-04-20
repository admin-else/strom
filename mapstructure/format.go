package mapstructure

import (
	"errors"
	"reflect"
)

var (
	ErrNotPointer      = errors.New("target must be a pointer")
	ErrNilPointer      = errors.New("target must not be a nil pointer")
	ErrNotStruct       = errors.New("target must be a pointer to a struct")
	ErrNotMap          = errors.New("data must be a map[string]any")
	ErrMissingRequired = errors.New("missing required field")
	ErrExtraKey        = errors.New("map contains key with no matching struct field")
)

type Format struct {
	Name         string
	requireAll   bool
	errorOnExtra bool
	decoders     map[reflect.Type]func(any, string) (any, error)
	encoders     map[reflect.Type]func(any, string) (any, error)
}

func NewFormat(name string, opts ...Option) *Format {
	f := &Format{
		Name:     name,
		decoders: make(map[reflect.Type]func(any, string) (any, error)),
		encoders: make(map[reflect.Type]func(any, string) (any, error)),
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

type Option func(*Format)

func WithRequireAll(require bool) Option {
	return func(f *Format) {
		f.requireAll = require
	}
}

func WithErrorOnExtra(errorOnExtra bool) Option {
	return func(f *Format) {
		f.errorOnExtra = errorOnExtra
	}
}

func WithTypeCodec(decodeFn, encodeFn any) Option {
	dt := reflect.TypeOf(decodeFn)
	et := reflect.TypeOf(encodeFn)

	if dt.Kind() != reflect.Func || dt.NumIn() != 2 || dt.NumOut() != 2 {
		panic("decode function must have signature func(any, string) (T, error)")
	}
	if et.Kind() != reflect.Func || et.NumIn() != 2 || et.NumOut() != 2 {
		panic("encode function must have signature func(T, string) (any, error)")
	}

	targetType := dt.Out(0)
	sourceType := et.In(0)

	return func(f *Format) {
		dv := reflect.ValueOf(decodeFn)
		f.decoders[targetType] = func(v any, tag string) (any, error) {
			results := dv.Call([]reflect.Value{reflect.ValueOf(v), reflect.ValueOf(tag)})
			err, _ := results[1].Interface().(error)
			return results[0].Interface(), err
		}

		ev := reflect.ValueOf(encodeFn)
		f.encoders[sourceType] = func(v any, tag string) (any, error) {
			results := ev.Call([]reflect.Value{reflect.ValueOf(v), reflect.ValueOf(tag)})
			err, _ := results[1].Interface().(error)
			return results[0].Interface(), err
		}
	}
}

func (f *Format) Decode(data any, target any) error {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return ErrNotPointer
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	m, ok := data.(map[string]any)
	if !ok {
		return ErrNotMap
	}
	return f.decode(m, v)
}

func (f *Format) Encode(source any) (map[string]any, error) {
	v := reflect.ValueOf(source)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil, ErrNilPointer
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}
	return f.encode(v)
}
