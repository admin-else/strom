package mapstructure

import "strings"

type tagOptions struct {
	Name      string
	OmitEmpty bool
	Required  bool
	Skip      bool
	Raw       string
}

func parseTag(raw string) tagOptions {
	if raw == "-" {
		return tagOptions{Skip: true, Raw: raw}
	}
	parts := strings.Split(raw, ",")
	opts := tagOptions{Name: parts[0], Raw: raw}
	for _, part := range parts[1:] {
		switch part {
		case "omitempty":
			opts.OmitEmpty = true
		case "required":
			opts.Required = true
		}
	}
	return opts
}
