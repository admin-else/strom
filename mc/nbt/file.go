package nbt

import (
	"compress/gzip"
	"io"
	"os"

	"github.com/admin-else/strom/mc/mapstructure"
)

var Format = mapstructure.NewFormat("nbt", mapstructure.WithRequireAll(), mapstructure.WithTrySnakeCase(), mapstructure.WithTryLowCase())

// ReadFile reads a gzipped NBT file from r and decodes it into data using the package-level Format.
func ReadFile(file io.Reader, data any) (err error) {
	return ReadFileAccounter(file, data, NbtAccounterUnlimitedHeap())
}

// ReadFileAccounter reads a gzipped NBT file from r and decodes it into data
// while charging the given NbtAccounter.
func ReadFileAccounter(file io.Reader, data any, accounter *NbtAccounter) (err error) {
	n, err := ReadUnstructuredFileAccounter(file, accounter)
	if err != nil {
		return
	}
	return Format.Decode(n.Value, data)
}

// ReadFileUncompressed reads an uncompressed NBT file from r and decodes it into
// data using the package-level Format (mirrors NbtIo.read + codec decode).
func ReadFileUncompressed(file io.Reader, data any) (err error) {
	return ReadFileUncompressedAccounter(file, data, NbtAccounterUnlimitedHeap())
}

// ReadFileUncompressedAccounter is ReadFileUncompressed with an NbtAccounter.
func ReadFileUncompressedAccounter(file io.Reader, data any, accounter *NbtAccounter) (err error) {
	n, err := ReadUncompressedFileAccounter(file, accounter)
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
	return ReadUnstructuredFileAccounter(file, NbtAccounterUnlimitedHeap())
}

// ReadUnstructuredFileAccounter is ReadUnstructuredFile with an NbtAccounter.
func ReadUnstructuredFileAccounter(file io.Reader, accounter *NbtAccounter) (n *Tag, err error) {
	r, err := gzip.NewReader(file)
	if err != nil {
		return
	}
	defer r.Close()
	n = &Tag{}
	err = n.DecodeAccounter(r, accounter)
	return
}

// ReadUnstructuredFilePath reads a gzipped NBT file from the given filesystem path and returns the root Tag.
func ReadUnstructuredFilePath(path string) (n *Tag, err error) {
	return ReadUnstructuredFilePathAccounter(path, NbtAccounterUnlimitedHeap())
}

// ReadUnstructuredFilePathAccounter is ReadUnstructuredFilePath with an NbtAccounter.
func ReadUnstructuredFilePathAccounter(path string, accounter *NbtAccounter) (n *Tag, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	return ReadUnstructuredFileAccounter(f, accounter)
}

// ReadUncompressedFile reads an uncompressed NBT file from r and returns the
// root Tag (mirrors NbtIo.read, which is not gzipped in this snapshot).
func ReadUncompressedFile(file io.Reader) (n *Tag, err error) {
	return ReadUncompressedFileAccounter(file, NbtAccounterUnlimitedHeap())
}

// ReadUncompressedFileAccounter is ReadUncompressedFile with an NbtAccounter.
func ReadUncompressedFileAccounter(file io.Reader, accounter *NbtAccounter) (n *Tag, err error) {
	n = &Tag{}
	err = n.DecodeAccounter(file, accounter)
	return
}

// ReadUncompressedFilePath reads an uncompressed NBT file from the given path
// and returns the root Tag.
func ReadUncompressedFilePath(path string) (n *Tag, err error) {
	return ReadUncompressedFilePathAccounter(path, NbtAccounterUnlimitedHeap())
}

// ReadUncompressedFilePathAccounter is ReadUncompressedFilePath with an NbtAccounter.
func ReadUncompressedFilePathAccounter(path string, accounter *NbtAccounter) (n *Tag, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	return ReadUncompressedFileAccounter(f, accounter)
}

// WriteUncompressedFile writes an NBT tag uncompressed to w (mirrors
// NbtIo.write(CompoundTag, DataOutput)).
func WriteUncompressedFile(file io.Writer, n *Tag) (err error) {
	return n.Encode(file)
}

// WriteUncompressedFilePath writes an NBT tag uncompressed to a filesystem path.
// The file is synced before close, mirroring NbtIo.SYNC_OUTPUT_OPTIONS.
func WriteUncompressedFilePath(path string, n *Tag) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	if err = WriteUncompressedFile(f, n); err != nil {
		f.Close()
		return
	}
	if err = f.Sync(); err != nil {
		f.Close()
		return
	}
	return f.Close()
}
