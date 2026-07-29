package mapstructure

var DefaultFormat = NewFormat("mapstructure", WithRequireAll(), WithTrySnakeCase(), WithTryLowCase())

// Decode is a stub placeholder. Use Format.Decode instead.
func Decode(data any, v any) error {
	return nil
}
