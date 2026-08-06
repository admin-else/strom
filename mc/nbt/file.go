package nbt

import (
	"compress/gzip"
	"io"
	"os"

	"git.anygate.cloud/anygatecloud/strom/mc/mapstructure"
)

var Format = mapstructure.NewFormat("nbt", mapstructure.WithRequireAll(), mapstructure.WithTrySnakeCase(), mapstructure.WithTryLowCase())

// ReadFile reads a gzipped NBT file from r and decodes it into data using the package-level Format.
func ReadFile(file io.Reader, data any) (err error) {
	n, err := ReadUnstructuredFile(file)
	if err != nil {
		return
	}
	return Format.Decode(n.Value, data)
}

// WriteUnstructuredFilePath writes an NBT tag as a gzipped file to the given filesystem path.
func WriteUnstructuredFilePath(path string, n *Tag) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	return WriteUnstructuredFile(f, n)
}

// WriteUnstructuredFile writes an NBT tag as a gzipped file to w.
func WriteUnstructuredFile(file io.Writer, n *Tag) (err error) {
	gzw := gzip.NewWriter(file)
	defer gzw.Close()
	err = n.Encode(gzw)
	return
}

// ReadUnstructuredFile reads a gzipped NBT file from r and returns the root Tag.
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

// ReadUnstructuredFilePath reads a gzipped NBT file from the given filesystem path and returns the root Tag.
func ReadUnstructuredFilePath(path string) (n *Tag, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	return ReadUnstructuredFile(f)
}
