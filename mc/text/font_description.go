package text

import "github.com/admin-else/strom/mc/resources"

// FontDescription mirrors net.minecraft.network.chat.FontDescription.
type FontDescription interface {
	fontDescription()
}

// FontDescriptionDEFAULT mirrors FontDescription.DEFAULT.
var FontDescriptionDEFAULT FontDescription = FontDescriptionResource{ID: resources.IdentifierWithDefaultNamespace("default")}

// FontDescriptionResource mirrors FontDescription.Resource.
type FontDescriptionResource struct {
	ID resources.Identifier
}

func (FontDescriptionResource) fontDescription() {}

// FontDescriptionAtlasSprite mirrors FontDescription.AtlasSprite.
type FontDescriptionAtlasSprite struct {
	AtlasID  resources.Identifier
	SpriteID resources.Identifier
}

func (FontDescriptionAtlasSprite) fontDescription() {}

// FontDescriptionPlayerSprite mirrors FontDescription.PlayerSprite. ResolvableProfile
// is not ported; Profile is left as an opaque value (see docs/QUESTIONS.md).
type FontDescriptionPlayerSprite struct {
	Profile any
	Hat     bool
}

func (FontDescriptionPlayerSprite) fontDescription() {}
