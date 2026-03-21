package nbt

import (
	"fmt"
	"io"
)

// PrintSNBTAny documented at https://minecraft.wiki/w/NBT_format#SNBT_format
func PrintSNBTAny(t any, w io.Writer) (err error) {
	switch t := t.(type) {
	case int8:
		_, err = fmt.Fprintf(w, "%vB", t)
	case bool:
		_, err = fmt.Fprintf(w, "%v", t)
	case int16:
		_, err = fmt.Fprintf(w, "%vS", t)
	case int32:
		_, err = fmt.Fprintf(w, "%vI", t)
	case int64:
		_, err = fmt.Fprintf(w, "%vL", t)
	case float32:
		_, err = fmt.Fprintf(w, "%vF", t)
	case float64:
		_, err = fmt.Fprintf(w, "%vD", t)
	case string:
		_, err = fmt.Fprintf(w, "%q", t)
	case []any:
		_, err = fmt.Fprintf(w, "[")
		for _, v := range t {
			err = PrintSNBTAny(v, w)
			if err != nil {
				return
			}
			_, err = fmt.Fprintf(w, ",")
		}
		_, err = fmt.Fprintf(w, "]")
	case map[string]any:
		_, err = fmt.Fprintf(w, "{")
		for k, v := range t {
			_, err = fmt.Fprintf(w, "%q:,", k)
			err = PrintSNBTAny(v, w)
			if err != nil {
				return
			}
			_, err = fmt.Fprintf(w, ",")
		}
		_, err = fmt.Fprintf(w, "}")
	case []int8:
		_, err = fmt.Fprintf(w, "[B;")
		if err != nil {
			return
		}
		for _, v := range t {
			err = PrintSNBTAny(v, w)
			if err != nil {
				return
			}
			_, err = fmt.Fprintf(w, ",")
		}
		_, err = fmt.Fprintf(w, "]")
	case []int32:
		_, err = fmt.Fprintf(w, "[I;")
		if err != nil {
			return
		}
		for _, v := range t {
			err = PrintSNBTAny(v, w)
			if err != nil {
				return
			}
			_, err = fmt.Fprintf(w, ",")
		}
		_, err = fmt.Fprintf(w, "]")
	case []int64:
		_, err = fmt.Fprintf(w, "[L;")
		if err != nil {
			return
		}
		for _, v := range t {
			err = PrintSNBTAny(v, w)
			if err != nil {
				return
			}
			_, err = fmt.Fprintf(w, ",")
		}
		_, err = fmt.Fprintf(w, "]")
	default:
		err = UnknownTagTypeErr
	}
	return
}
