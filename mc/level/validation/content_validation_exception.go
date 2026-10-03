package validation

import "strings"

// ContentValidationException mirrors net.minecraft.world.level.validation.ContentValidationException.
type ContentValidationException struct {
	directory string
	entries   []ForbiddenSymlinkInfo
}

// NewContentValidationException mirrors the ContentValidationException(Path, List<ForbiddenSymlinkInfo>) constructor.
func NewContentValidationException(directory string, entries []ForbiddenSymlinkInfo) *ContentValidationException {
	return &ContentValidationException{directory: directory, entries: entries}
}

// Error mirrors ContentValidationException.getMessage().
func (e *ContentValidationException) Error() string {
	return ContentValidationExceptionGetMessage(e.directory, e.entries)
}

// ContentValidationExceptionGetMessage mirrors the static ContentValidationException.getMessage(Path, List<ForbiddenSymlinkInfo>).
func ContentValidationExceptionGetMessage(directory string, entries []ForbiddenSymlinkInfo) string {
	parts := make([]string, len(entries))
	for i, entry := range entries {
		parts[i] = entry.Link() + "->" + entry.Target()
	}
	return "Failed to validate '" + directory + "'. Found forbidden symlinks: " + strings.Join(parts, ", ")
}
