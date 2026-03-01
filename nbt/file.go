package nbt

import (
	"compress/gzip"
	"io"
)

func WriteFile(file io.Writer, n Tag) (err error) {
	gzw := gzip.NewWriter(file)
	if err != nil {
		return
	}
	defer gzw.Close()
	err = n.Encode(gzw)
	return
}

func ReadFile(file io.Reader) (err error, n Tag) {
	r, err := gzip.NewReader(file)
	if err != nil {
		return
	}
	defer r.Close()
	t := &Tag{}
	err = t.Decode(r)
	return
}
