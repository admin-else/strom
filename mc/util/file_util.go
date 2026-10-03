package util

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	reservedWindowsFilenames = regexp.MustCompile(`^(?i)(?:.*\.|(?:COM|CLOCK\$|CON|PRN|AUX|NUL|COM[1-9]|LPT[1-9])(?:\..*)?)$`)
	strictPathSegmentCheck   = regexp.MustCompile(`^[-._a-z0-9]+$`)
)

// DecomposePath mirrors FileUtil.decomposePath(String). Java returns a
// DataResult<List<String>>; Go returns the list or a plain error.
func DecomposePath(path string) (result []string, err error) {
	segmentEnd := indexByte(path, '/')
	if segmentEnd == -1 {
		switch path {
		case "", ".", "..":
			return nil, fmt.Errorf("Invalid path '%s'", path)
		default:
			if !containsAllowedCharactersOnly(path) {
				return nil, fmt.Errorf("Invalid path '%s'", path)
			}
			return []string{path}, nil
		}
	}

	result = make([]string, 0, 4)
	segmentStart := 0
	lastSegment := false
	for {
		segment := path[segmentStart:segmentEnd]
		switch segment {
		case "", ".", "..":
			return nil, fmt.Errorf("Invalid segment '%s' in path '%s'", segment, path)
		}
		if !containsAllowedCharactersOnly(segment) {
			return nil, fmt.Errorf("Invalid segment '%s' in path '%s'", segment, path)
		}
		result = append(result, segment)
		if lastSegment {
			return result, nil
		}
		segmentStart = segmentEnd + 1
		next := indexByteFrom(path, '/', segmentStart)
		if next == -1 {
			segmentEnd = len(path)
			lastSegment = true
		} else {
			segmentEnd = next
		}
	}
}

// ResolvePath mirrors FileUtil.resolvePath(Path, List<String>).
func ResolvePath(root string, segments []string) string {
	switch len(segments) {
	case 0:
		return root
	case 1:
		return filepath.Join(root, segments[0])
	default:
		return filepath.Join(root, filepath.Join(segments...))
	}
}

// IsValidPathSegment mirrors FileUtil.isValidPathSegment(String).
func IsValidPathSegment(segment string) bool {
	return segment != ".." && segment != "." && containsAllowedCharactersOnly(segment)
}

// ValidatePath mirrors FileUtil.validatePath(String...). Java throws
// IllegalArgumentException; Go returns a plain error.
func ValidatePath(path ...string) error {
	if len(path) == 0 {
		return fmt.Errorf("Path must have at least one element")
	}
	for _, segment := range path {
		if !IsValidPathSegment(segment) {
			return fmt.Errorf("Illegal segment %s in path [%s]", segment, strings.Join(path, ", "))
		}
	}
	return nil
}

// IsPathPartPortable mirrors FileUtil.isPathPartPortable(String).
func IsPathPartPortable(name string) bool {
	return !reservedWindowsFilenames.MatchString(name)
}

// CreateDirectoriesSafe mirrors FileUtil.createDirectoriesSafe(Path): create
// the directory tree of dir.toRealPath() when dir already exists, otherwise of
// dir itself.
func CreateDirectoriesSafe(dir string) (err error) {
	target := dir
	if _, statErr := os.Lstat(dir); statErr == nil {
		if real, realErr := filepath.EvalSymlinks(dir); realErr == nil {
			target = real
		}
	}
	return os.MkdirAll(target, 0o755)
}

func containsAllowedCharactersOnly(segment string) bool {
	return strictPathSegmentCheck.MatchString(segment)
}

func indexByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}

func indexByteFrom(s string, b byte, from int) int {
	for i := from; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}
