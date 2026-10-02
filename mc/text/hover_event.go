package text

import (
	"fmt"
	"reflect"

	"github.com/google/uuid"
)

// HoverEvent mirrors net.minecraft.network.chat.HoverEvent.
type HoverEvent interface {
	Action() HoverEventAction
}

// HoverEventAction mirrors HoverEvent.Action.
type HoverEventAction struct {
	name            string
	allowFromServer bool
}

var (
	HoverEventActionSHOW_TEXT   = HoverEventAction{"show_text", true}
	HoverEventActionSHOW_ITEM   = HoverEventAction{"show_item", true}
	HoverEventActionSHOW_ENTITY = HoverEventAction{"show_entity", true}
)

// GetSerializedName mirrors HoverEvent.Action.getSerializedName().
func (a HoverEventAction) GetSerializedName() string { return a.name }

// IsAllowedFromServer mirrors HoverEvent.Action.isAllowedFromServer().
func (a HoverEventAction) IsAllowedFromServer() bool { return a.allowFromServer }

// String mirrors HoverEvent.Action.toString().
func (a HoverEventAction) String() string { return "<action " + a.name + ">" }

// HoverEventActionFilterForSerialization mirrors HoverEvent.Action.filterForSerialization.
func HoverEventActionFilterForSerialization(action HoverEventAction) (HoverEventAction, error) {
	if !action.IsAllowedFromServer() {
		return HoverEventAction{}, fmt.Errorf("Action not allowed: %s", action.String())
	}
	return action, nil
}

// HoverEventEntityTooltipInfo mirrors HoverEvent.EntityTooltipInfo. EntityType is not
// ported; Type stays opaque and getTooltipLines() is deferred (see docs/QUESTIONS.md).
type HoverEventEntityTooltipInfo struct {
	Type any
	UUID uuid.UUID
	Name *Component
}

// NewHoverEventEntityTooltipInfo mirrors the EntityTooltipInfo(EntityType, UUID, Component)
// constructor (Optional.ofNullable(name)).
func NewHoverEventEntityTooltipInfo(entityType any, id uuid.UUID, name *Component) HoverEventEntityTooltipInfo {
	return HoverEventEntityTooltipInfo{Type: entityType, UUID: id, Name: name}
}

// Equals mirrors EntityTooltipInfo.equals.
func (e HoverEventEntityTooltipInfo) Equals(o HoverEventEntityTooltipInfo) bool {
	if !reflect.DeepEqual(e.Type, o.Type) {
		return false
	}
	if e.UUID != o.UUID {
		return false
	}
	return optionalComponentEquals(e.Name, o.Name)
}

// HashCode mirrors EntityTooltipInfo.hashCode.
func (e HoverEventEntityTooltipInfo) HashCode() int {
	h := stringHash(fmt.Sprintf("%v", e.Type))
	h = 31*h + stringHash(e.UUID.String())
	h = 31*h + optionalComponentHash(e.Name)
	return h
}

// HoverEventShowText mirrors HoverEvent.ShowText.
type HoverEventShowText struct {
	Value Component
}

func (HoverEventShowText) Action() HoverEventAction { return HoverEventActionSHOW_TEXT }

// HoverEventShowItem mirrors HoverEvent.ShowItem. ItemStackTemplate is not ported;
// Item stays opaque (see docs/QUESTIONS.md).
type HoverEventShowItem struct {
	Item any
}

func (HoverEventShowItem) Action() HoverEventAction { return HoverEventActionSHOW_ITEM }

// HoverEventShowEntity mirrors HoverEvent.ShowEntity.
type HoverEventShowEntity struct {
	Entity HoverEventEntityTooltipInfo
}

func (HoverEventShowEntity) Action() HoverEventAction { return HoverEventActionSHOW_ENTITY }

// String mirrors the HoverEvent.ShowText record toString.
func (e HoverEventShowText) String() string { return fmt.Sprintf("ShowText[value=%v]", e.Value) }

// String mirrors the HoverEvent.ShowItem record toString.
func (e HoverEventShowItem) String() string { return fmt.Sprintf("ShowItem[item=%v]", e.Item) }

// String mirrors the HoverEvent.ShowEntity record toString.
func (e HoverEventShowEntity) String() string { return fmt.Sprintf("ShowEntity[entity=%v]", e.Entity) }
