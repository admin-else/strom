package text

// KeybindContents mirrors net.minecraft.network.chat.contents.KeybindContents.
type KeybindContents struct {
	BaseComponentContents
	name         string
	nameResolver func() Component
}

// NewKeybindContents mirrors the KeybindContents constructor.
func NewKeybindContents(name string) *KeybindContents {
	return &KeybindContents{name: name}
}

func (k *KeybindContents) getNestedComponent() Component {
	if k.nameResolver == nil {
		k.nameResolver = keybindKeyResolver(k.name)
	}
	return k.nameResolver()
}

// VisitContent mirrors KeybindContents.visit(ContentConsumer).
func (k *KeybindContents) VisitContent(output ContentConsumer) (any, bool) {
	return k.getNestedComponent().VisitContent(output)
}

// VisitStyled mirrors KeybindContents.visit(StyledContentConsumer, Style).
func (k *KeybindContents) VisitStyled(output StyledContentConsumer, currentStyle Style) (any, bool) {
	return k.getNestedComponent().VisitStyled(output, currentStyle)
}

// Resolve mirrors the ComponentContents default.
func (k *KeybindContents) Resolve(context *ResolutionContext, recursionDepth int) (*MutableComponent, error) {
	return componentContentsSelfResolve(k), nil
}

// Equals mirrors KeybindContents.equals.
func (k *KeybindContents) Equals(o *KeybindContents) bool {
	if k == o {
		return true
	}
	return o != nil && k.name == o.name
}

// HashCode mirrors KeybindContents.hashCode.
func (k *KeybindContents) HashCode() int { return stringHash(k.name) }

func (k *KeybindContents) String() string { return "keybind{" + k.name + "}" }

// GetName mirrors KeybindContents.getName().
func (k *KeybindContents) GetName() string { return k.name }
