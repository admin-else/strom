package text

import (
	"fmt"

	"github.com/admin-else/strom/mc/nbt"
	"github.com/admin-else/strom/mc/resources"
)

// ClickEvent mirrors net.minecraft.network.chat.ClickEvent.
type ClickEvent interface {
	Action() ClickEventAction
}

// ClickEventAction mirrors ClickEvent.Action.
type ClickEventAction struct {
	name            string
	serializedName  string
	allowFromServer bool
}

var (
	ClickEventActionOPEN_URL          = ClickEventAction{"OPEN_URL", "open_url", true}
	ClickEventActionOPEN_FILE         = ClickEventAction{"OPEN_FILE", "open_file", false}
	ClickEventActionRUN_COMMAND       = ClickEventAction{"RUN_COMMAND", "run_command", true}
	ClickEventActionSUGGEST_COMMAND   = ClickEventAction{"SUGGEST_COMMAND", "suggest_command", true}
	ClickEventActionSHOW_DIALOG       = ClickEventAction{"SHOW_DIALOG", "show_dialog", true}
	ClickEventActionCHANGE_PAGE       = ClickEventAction{"CHANGE_PAGE", "change_page", true}
	ClickEventActionCOPY_TO_CLIPBOARD = ClickEventAction{"COPY_TO_CLIPBOARD", "copy_to_clipboard", true}
	ClickEventActionCUSTOM            = ClickEventAction{"CUSTOM", "custom", true}
)

// GetSerializedName mirrors ClickEvent.Action.getSerializedName().
func (a ClickEventAction) GetSerializedName() string { return a.serializedName }

// String mirrors the ClickEvent.Action enum toString (name()).
func (a ClickEventAction) String() string { return a.name }

// IsAllowedFromServer mirrors ClickEvent.Action.isAllowedFromServer().
func (a ClickEventAction) IsAllowedFromServer() bool { return a.allowFromServer }

// ClickEventActionFilterForSerialization mirrors ClickEvent.Action.filterForSerialization.
func ClickEventActionFilterForSerialization(action ClickEventAction) (ClickEventAction, error) {
	if !action.IsAllowedFromServer() {
		return ClickEventAction{}, fmt.Errorf("Click event type not allowed: %s", action.String())
	}
	return action, nil
}

// ClickEventChangePage mirrors ClickEvent.ChangePage.
type ClickEventChangePage struct {
	Page int
}

func (ClickEventChangePage) Action() ClickEventAction { return ClickEventActionCHANGE_PAGE }

// ClickEventCopyToClipboard mirrors ClickEvent.CopyToClipboard.
type ClickEventCopyToClipboard struct {
	Value string
}

func (ClickEventCopyToClipboard) Action() ClickEventAction { return ClickEventActionCOPY_TO_CLIPBOARD }

// ClickEventCustom mirrors ClickEvent.Custom.
type ClickEventCustom struct {
	ID      resources.Identifier
	Payload *nbt.Tag
}

func (ClickEventCustom) Action() ClickEventAction { return ClickEventActionCUSTOM }

// ClickEventOpenFile mirrors ClickEvent.OpenFile.
type ClickEventOpenFile struct {
	Path string
}

func (ClickEventOpenFile) Action() ClickEventAction { return ClickEventActionOPEN_FILE }

// ClickEventOpenUrl mirrors ClickEvent.OpenUrl.
type ClickEventOpenUrl struct {
	URI string
}

func (ClickEventOpenUrl) Action() ClickEventAction { return ClickEventActionOPEN_URL }

// ClickEventRunCommand mirrors ClickEvent.RunCommand.
type ClickEventRunCommand struct {
	Command string
}

func (ClickEventRunCommand) Action() ClickEventAction { return ClickEventActionRUN_COMMAND }

// ClickEventShowDialog mirrors ClickEvent.ShowDialog. Holder<Dialog> is not ported;
// Dialog is an opaque value (see docs/QUESTIONS.md).
type ClickEventShowDialog struct {
	Dialog any
}

func (ClickEventShowDialog) Action() ClickEventAction { return ClickEventActionSHOW_DIALOG }

// ClickEventSuggestCommand mirrors ClickEvent.SuggestCommand.
type ClickEventSuggestCommand struct {
	Command string
}

func (ClickEventSuggestCommand) Action() ClickEventAction { return ClickEventActionSUGGEST_COMMAND }

// String mirrors the ClickEvent.ChangePage record toString.
func (e ClickEventChangePage) String() string { return fmt.Sprintf("ChangePage[page=%d]", e.Page) }

// String mirrors the ClickEvent.CopyToClipboard record toString.
func (e ClickEventCopyToClipboard) String() string {
	return "CopyToClipboard[value=" + e.Value + "]"
}

// String mirrors the ClickEvent.Custom record toString.
func (e ClickEventCustom) String() string {
	return fmt.Sprintf("Custom[id=%s, payload=%v]", e.ID.String(), e.Payload)
}

// String mirrors the ClickEvent.OpenFile record toString.
func (e ClickEventOpenFile) String() string { return "OpenFile[path=" + e.Path + "]" }

// String mirrors the ClickEvent.OpenUrl record toString.
func (e ClickEventOpenUrl) String() string { return "OpenUrl[uri=" + e.URI + "]" }

// String mirrors the ClickEvent.RunCommand record toString.
func (e ClickEventRunCommand) String() string { return "RunCommand[command=" + e.Command + "]" }

// String mirrors the ClickEvent.ShowDialog record toString.
func (e ClickEventShowDialog) String() string { return fmt.Sprintf("ShowDialog[dialog=%v]", e.Dialog) }

// String mirrors the ClickEvent.SuggestCommand record toString.
func (e ClickEventSuggestCommand) String() string {
	return "SuggestCommand[command=" + e.Command + "]"
}
