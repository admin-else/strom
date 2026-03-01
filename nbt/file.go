package nbt

import (
	"compress/gzip"
	"io"
)

func WriteFile(file io.Writer, n Tag) (err error) {
	gzw := gzip.NewWriter(file)
	defer gzw.Close()
	err = n.Encode(gzw)
	return
}

func ReadFile(file io.Reader) (n *Tag, err error) {
	r, err := gzip.NewReader(file)
	if err != nil {
		return
	}
	defer r.Close()
	n = &Tag{}
	err = n.Decode(r)
	return
}
