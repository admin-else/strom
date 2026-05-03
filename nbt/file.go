package nbt

import (
	"compress/gzip"
	"io"
	"os"

	"github.com/admin-else/strom/mapstructure"
)

var Format = mapstructure.NewFormat("nbt", mapstructure.WithRequireAll(), mapstructure.WithTrySnakeCase(), mapstructure.WithTryLowCase())

func ReadFile(file io.Reader, data any) (err error) {
	n, err := ReadUnstructuredFile(file)
	if err != nil {
		return
	}
	return Format.Decode(n.Value, data)
}

func WriteUnstructuredFilePath(path string, n *Tag) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	return WriteUnstructuredFile(f, n)
}

func WriteUnstructuredFile(file io.Writer, n *Tag) (err error) {
	gzw := gzip.NewWriter(file)
	defer gzw.Close()
	err = n.Encode(gzw)
	return
}

func ReadUnstructuredFile(file io.Reader) (n *Tag, err error) {
	r, err := gzip.NewReader(file)
	if err != nil {
		return
	}
	defer r.Close()
	n = &Tag{}
	err = n.Decode(r)
	return
}

func ReadUnstructuredFilePath(path string) (n *Tag, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	return ReadUnstructuredFile(f)
}
