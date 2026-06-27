package mapstructure

import (
	"fmt"
	"reflect"
)

func (f *Format) encode(v reflect.Value) (map[string]any, error) {
	t := v.Type()
	result := map[string]any{}

	for i := range t.NumField() {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		tag := parseTag(field.Tag.Get(f.Name))
		if tag.Skip {
			continue
		}

		key := tag.Name
		if key == "" {
			key = field.Name
		}

		fieldVal := v.Field(i)
		if tag.OmitEmpty && fieldVal.IsZero() {
			continue
		}

		encoded, err := f.encodeField(fieldVal, tag.Raw)
		if err != nil {
			return nil, fmt.Errorf("field %s: %w", key, err)
		}
		result[key] = encoded
	}

	return result, nil
}

func (f *Format) encodeField(v reflect.Value, tagStr string) (any, error) {
	if codec, ok := f.encoders[v.Type()]; ok {
		return codec(v.Interface(), tagStr)
	}

	switch v.Kind() {
	case reflect.Struct:
		return f.encode(v)
	case reflect.Slice:
		if v.IsNil() {
			return []any{}, nil
		}
		slice := make([]any, v.Len())
		for i := range v.Len() {
			encoded, err := f.encodeField(v.Index(i), tagStr)
			if err != nil {
				return nil, err
			}
			slice[i] = encoded
		}
		return slice, nil
	case reflect.Map:
		if v.IsNil() {
			return map[string]any{}, nil
		}
		m := make(map[string]any, v.Len())
		iter := v.MapRange()
		for iter.Next() {
			key := iter.Key().String()
			encoded, err := f.encodeField(iter.Value(), tagStr)
			if err != nil {
				return nil, err
			}
			m[key] = encoded
		}
		return m, nil
	case reflect.Pointer:
		if v.IsNil() {
			return nil, nil
		}
		return f.encodeField(v.Elem(), tagStr)
	case reflect.Bool:
		if v.Bool() {
			return int8(1), nil
		}
		return int8(0), nil
	default:
		return v.Interface(), nil
	}
}
