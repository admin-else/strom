package mapstructure

var DefaultFormat = NewFormat("mapstructure", WithRequireAll(), WithTrySnakeCase(), WithTryLowCase())

func Decode(data any, v any) error {
	return nil
}
