package mapstructure

import (
	"fmt"
	"reflect"

	"github.com/admin-else/strom/util"
)

func (f *Format) decode(m map[string]any, target reflect.Value) error {
	t := target.Type()
	usedKeys := map[string]bool{}

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

		val, exists := m[key]
		if !exists && f.trySnakeCase {
			val, exists = m[util.SnakeCase(key)]
		}
		if !exists && f.tryLowCase {
			val, exists = m[util.FirstLetterLower(key)]
		}
		if !exists {
			if tag.Required || (f.requireAll && !tag.OmitEmpty) {
				return fmt.Errorf("%w: %s", ErrMissingRequired, key)
			}
			continue
		}

		usedKeys[key] = true

		fieldVal := target.Field(i)
		if err := f.setField(fieldVal, val, tag); err != nil {
			return fmt.Errorf("field %s: %w", key, err)
		}
	}

	if f.errorOnExtra {
		for key := range m {
			if !usedKeys[key] {
				return fmt.Errorf("%w: %s", ErrExtraKey, key)
			}
		}
	}

	return nil
}

func (f *Format) setField(target reflect.Value, value any, opts tagOptions) error {
	if value == nil {
		if opts.OmitEmpty {
			return nil
		}
		target.Set(reflect.Zero(target.Type()))
		return nil
	}

	targetType := target.Type()

	if codec, ok := f.decoders[targetType]; ok {
		result, err := codec(value, opts.Raw)
		if err != nil {
			return err
		}
		target.Set(reflect.ValueOf(result))
		return nil
	}

	src := reflect.ValueOf(value)

	if src.Type().AssignableTo(targetType) {
		target.Set(src)
		return nil
	}

	if src.Type().ConvertibleTo(targetType) {
		target.Set(src.Convert(targetType))
		return nil
	}

	if targetType.Kind() == reflect.Struct {
		nested, ok := value.(map[string]any)
		if !ok {
			return fmt.Errorf("cannot map %T to struct %s", value, targetType.Name())
		}
		nestedTarget := reflect.New(targetType).Elem()
		if err := f.decode(nested, nestedTarget); err != nil {
			return err
		}
		target.Set(nestedTarget)
		return nil
	}

	if targetType.Kind() == reflect.Slice {
		return f.setSlice(target, value)
	}

	if targetType.Kind() == reflect.Map {
		return f.setMap(target, value)
	}

	if targetType.Kind() == reflect.Bool {
		return f.setBool(target, value)
	}

	// this fixes this edge case
	//     level_test.go:19: field Data: field Player: field Pos: cannot set [3]float64 from []interface {}
	if targetType.Kind() == reflect.Array && src.Kind() == reflect.Slice {
		if target.Len() != src.Len() {
			return fmt.Errorf("cannot set %s from %T because lengths do not match", targetType, value)
		}
		for i := 0; i < src.Len(); i++ {
			if src.Index(i).Type().AssignableTo(targetType.Elem()) {
				target.Index(i).Set(src.Index(i))
				continue
			}
			if src.Index(i).Elem().Type().AssignableTo(targetType.Elem()) {
				target.Index(i).Set(src.Index(i).Elem())
				continue
			}
			return fmt.Errorf("cannot set %s from %T because element %d cannot be converted to %s", targetType, value, i, targetType.Elem())
		}
		return nil
	}

	return fmt.Errorf("cannot set %s from %T", targetType, value)
}

func (f *Format) setBool(target reflect.Value, value any) error {
	switch v := value.(type) {
	case int8:
		target.SetBool(v != 0)
	case int16:
		target.SetBool(v != 0)
	case int32:
		target.SetBool(v != 0)
	case int64:
		target.SetBool(v != 0)
	default:
		return fmt.Errorf("cannot convert %T to bool", value)
	}
	return nil
}

func (f *Format) setSlice(target reflect.Value, value any) error {
	src, ok := value.([]any)
	if !ok {
		return fmt.Errorf("cannot map %T to slice", value)
	}

	elemType := target.Type().Elem()
	slice := reflect.MakeSlice(target.Type(), 0, len(src))

	for _, item := range src {
		elem := reflect.New(elemType).Elem()
		if err := f.setField(elem, item, tagOptions{}); err != nil {
			return err
		}
		slice = reflect.Append(slice, elem)
	}

	target.Set(slice)
	return nil
}

func (f *Format) setMap(target reflect.Value, value any) error {
	src, ok := value.(map[string]any)
	if !ok {
		return fmt.Errorf("cannot map %T to map", value)
	}

	mapType := target.Type()
	keyType := mapType.Key()
	elemType := mapType.Elem()

	result := reflect.MakeMap(mapType)

	for k, v := range src {
		keyVal := reflect.ValueOf(k)
		if !keyVal.Type().AssignableTo(keyType) {
			keyVal = keyVal.Convert(keyType)
		}

		elemVal := reflect.New(elemType).Elem()
		if err := f.setField(elemVal, v, tagOptions{}); err != nil {
			return err
		}

		result.SetMapIndex(keyVal, elemVal)
	}

	target.Set(result)
	return nil
}
