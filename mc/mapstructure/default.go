package mapstructure

var DefaultFormat = NewFormat("mapstructure", WithTrySnakeCase(), WithTryLowCase(), WithTryCaseInsensitive())

// Decode populates v from data using DefaultFormat.
func Decode(data any, v any) error {
	return DefaultFormat.Decode(data, v)
}
