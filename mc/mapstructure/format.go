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
	Name               string
	requireAll         bool
	errorOnExtra       bool
	trySnakeCase       bool
	tryLowCase         bool
	tryCaseInsensitive bool

	decoders map[reflect.Type]func(any, string) (any, error)
	encoders map[reflect.Type]func(any, string) (any, error)
}

// NewFormat creates a new Format with the given name and options.
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

// WithRequireAll returns an Option that requires all struct fields to be present in the map.
func WithRequireAll() Option {
	return func(f *Format) {
		f.requireAll = true
	}
}

// WithErrorOnExtra returns an Option that causes an error when the map contains unmatched keys.
func WithErrorOnExtra() Option {
	return func(f *Format) {
		f.errorOnExtra = true
	}
}

// WithTrySnakeCase returns an Option that tries to match struct fields using snake_case names.
func WithTrySnakeCase() Option {
	return func(f *Format) {
		f.trySnakeCase = true
	}
}

// WithTryLowCase returns an Option that tries to match struct fields using lowercase names.
func WithTryLowCase() Option {
	return func(f *Format) {
		f.tryLowCase = true
	}
}

// WithTryCaseInsensitive returns an Option that tries to match struct fields
// ignoring case, mirroring go-viper/mapstructure's default name matching.
func WithTryCaseInsensitive() Option {
	return func(f *Format) {
		f.tryCaseInsensitive = true
	}
}

// WithTypeCodec returns an Option that registers custom decode/encode functions for a type.
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

// Decode populates target (a pointer to a struct, slice or map) from data.
func (f *Format) Decode(data any, target any) error {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return ErrNotPointer
	}
	v = v.Elem()
	switch v.Kind() {
	case reflect.Struct:
		m, ok := data.(map[string]any)
		if !ok {
			return ErrNotMap
		}
		return f.decode(m, v)
	case reflect.Slice:
		return f.setSlice(v, data)
	case reflect.Map:
		return f.setMap(v, data)
	default:
		return ErrNotStruct
	}
}

// Encode converts source (a struct or pointer to a struct) into a map[string]any.
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
