package proto

import (
	"reflect"
)

func SmartConvertibleTo(from, to reflect.Type) bool {
	switch from.Kind() {
	case reflect.Struct:
		if from.NumField() != to.NumField() {
			return false
		}
		for i := range from.NumField() {
			if !SmartConvertibleTo(from.Field(i).Type, to.Field(i).Type) {
				return false
			}
		}
		return true
	case reflect.Ptr:
		return SmartConvertibleTo(from.Elem(), to.Elem())
	default:
		return from.ConvertibleTo(to)
	}
}

func SmartConvert(from reflect.Value, to reflect.Type) (ret reflect.Value) {
	switch from.Kind() {
	case reflect.Struct:
		ret = reflect.New(to).Elem()
		for i := range from.NumField() {
			ret.Field(i).Set(SmartConvert(from.Field(i), to.Field(i).Type))
		}
		return ret
	case reflect.Ptr:
		if from.IsNil() {
			return reflect.Zero(to)
		}
		return SmartConvert(from.Elem(), to.Elem()).Addr()
	default:
		return from.Convert(to)
	}
}
