// Package validation mirrors net.minecraft.world.level.validation.
package validation

// ForbiddenSymlinkInfo mirrors net.minecraft.world.level.validation.ForbiddenSymlinkInfo.
type ForbiddenSymlinkInfo struct {
	link   string
	target string
}

// NewForbiddenSymlinkInfo mirrors the ForbiddenSymlinkInfo(Path, Path) record constructor.
func NewForbiddenSymlinkInfo(link string, target string) ForbiddenSymlinkInfo {
	return ForbiddenSymlinkInfo{link: link, target: target}
}

// Link mirrors ForbiddenSymlinkInfo.link().
func (f ForbiddenSymlinkInfo) Link() string { return f.link }

// Target mirrors ForbiddenSymlinkInfo.target().
func (f ForbiddenSymlinkInfo) Target() string { return f.target }
