package validation

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// DirectoryValidator mirrors net.minecraft.world.level.validation.DirectoryValidator.
//
// Java takes a java.nio.file.PathMatcher for the symlink target allow-list. Go
// has no equivalent PathMatcher in the standard library, so the matcher is a
// plain func over the link target path.
type DirectoryValidator struct {
	symlinkTargetAllowList func(string) bool
}

// NewDirectoryValidator mirrors the DirectoryValidator(PathMatcher) constructor.
func NewDirectoryValidator(symlinkTargetAllowList func(string) bool) *DirectoryValidator {
	return &DirectoryValidator{symlinkTargetAllowList: symlinkTargetAllowList}
}

// ValidateSymlink mirrors DirectoryValidator.validateSymlink(Path, List<ForbiddenSymlinkInfo>).
func (d *DirectoryValidator) ValidateSymlink(path string, issues *[]ForbiddenSymlinkInfo) (err error) {
	target, err := os.Readlink(path)
	if err != nil {
		return err
	}
	if !d.symlinkTargetAllowList(target) {
		*issues = append(*issues, NewForbiddenSymlinkInfo(path, target))
	}
	return nil
}

// ValidateSymlinkList mirrors DirectoryValidator.validateSymlink(Path).
func (d *DirectoryValidator) ValidateSymlinkList(path string) (result []ForbiddenSymlinkInfo, err error) {
	if err = d.ValidateSymlink(path, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// ValidateDirectory mirrors DirectoryValidator.validateDirectory(Path, boolean).
func (d *DirectoryValidator) ValidateDirectory(directory string, allowTopSymlink bool) (issues []ForbiddenSymlinkInfo, err error) {
	attributes, err := os.Lstat(directory)
	if err != nil {
		if os.IsNotExist(err) {
			return issues, nil
		}
		return nil, err
	}

	if attributes.Mode().IsRegular() {
		return nil, fmt.Errorf("Path %s is not a directory", directory)
	}

	if attributes.Mode()&fs.ModeSymlink != 0 {
		if !allowTopSymlink {
			if err = d.ValidateSymlink(directory, &issues); err != nil {
				return nil, err
			}
			return issues, nil
		}

		if directory, err = os.Readlink(directory); err != nil {
			return nil, err
		}
	}

	if err = d.ValidateKnownDirectory(directory, &issues); err != nil {
		return nil, err
	}
	return issues, nil
}

// ValidateKnownDirectory mirrors DirectoryValidator.validateKnownDirectory(Path, List<ForbiddenSymlinkInfo>).
func (d *DirectoryValidator) ValidateKnownDirectory(directory string, issues *[]ForbiddenSymlinkInfo) (err error) {
	return filepath.WalkDir(directory, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		info, infoErr := entry.Info()
		if infoErr != nil {
			return infoErr
		}
		if info.Mode()&fs.ModeSymlink != 0 {
			if symlinkErr := d.ValidateSymlink(path, issues); symlinkErr != nil {
				return symlinkErr
			}
		}
		return nil
	})
}
