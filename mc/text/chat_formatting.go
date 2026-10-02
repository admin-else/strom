package text

import "regexp"

// ChatFormatting mirrors net.minecraft.ChatFormatting.
type ChatFormatting struct {
	name     string
	code     rune
	toString string
}

// ChatFormattingValues mirrors the implicit Enum.values().
func ChatFormattingValues() []ChatFormatting {
	return []ChatFormatting{
		ChatFormattingBLACK,
		ChatFormattingDARK_BLUE,
		ChatFormattingDARK_GREEN,
		ChatFormattingDARK_AQUA,
		ChatFormattingDARK_RED,
		ChatFormattingDARK_PURPLE,
		ChatFormattingGOLD,
		ChatFormattingGRAY,
		ChatFormattingDARK_GRAY,
		ChatFormattingBLUE,
		ChatFormattingGREEN,
		ChatFormattingAQUA,
		ChatFormattingRED,
		ChatFormattingLIGHT_PURPLE,
		ChatFormattingYELLOW,
		ChatFormattingWHITE,
		ChatFormattingOBFUSCATED,
		ChatFormattingBOLD,
		ChatFormattingSTRIKETHROUGH,
		ChatFormattingUNDERLINE,
		ChatFormattingITALIC,
		ChatFormattingRESET,
	}
}

var (
	ChatFormattingBLACK         = ChatFormatting{"BLACK", '0', "§0"}
	ChatFormattingDARK_BLUE     = ChatFormatting{"DARK_BLUE", '1', "§1"}
	ChatFormattingDARK_GREEN    = ChatFormatting{"DARK_GREEN", '2', "§2"}
	ChatFormattingDARK_AQUA     = ChatFormatting{"DARK_AQUA", '3', "§3"}
	ChatFormattingDARK_RED      = ChatFormatting{"DARK_RED", '4', "§4"}
	ChatFormattingDARK_PURPLE   = ChatFormatting{"DARK_PURPLE", '5', "§5"}
	ChatFormattingGOLD          = ChatFormatting{"GOLD", '6', "§6"}
	ChatFormattingGRAY          = ChatFormatting{"GRAY", '7', "§7"}
	ChatFormattingDARK_GRAY     = ChatFormatting{"DARK_GRAY", '8', "§8"}
	ChatFormattingBLUE          = ChatFormatting{"BLUE", '9', "§9"}
	ChatFormattingGREEN         = ChatFormatting{"GREEN", 'a', "§a"}
	ChatFormattingAQUA          = ChatFormatting{"AQUA", 'b', "§b"}
	ChatFormattingRED           = ChatFormatting{"RED", 'c', "§c"}
	ChatFormattingLIGHT_PURPLE  = ChatFormatting{"LIGHT_PURPLE", 'd', "§d"}
	ChatFormattingYELLOW        = ChatFormatting{"YELLOW", 'e', "§e"}
	ChatFormattingWHITE         = ChatFormatting{"WHITE", 'f', "§f"}
	ChatFormattingOBFUSCATED    = ChatFormatting{"OBFUSCATED", 'k', "§k"}
	ChatFormattingBOLD          = ChatFormatting{"BOLD", 'l', "§l"}
	ChatFormattingSTRIKETHROUGH = ChatFormatting{"STRIKETHROUGH", 'm', "§m"}
	ChatFormattingUNDERLINE     = ChatFormatting{"UNDERLINE", 'n', "§n"}
	ChatFormattingITALIC        = ChatFormatting{"ITALIC", 'o', "§o"}
	ChatFormattingRESET         = ChatFormatting{"RESET", 'r', "§r"}
)

// ChatFormattingPREFIX_CODE mirrors ChatFormatting.PREFIX_CODE.
const ChatFormattingPREFIX_CODE = '§'

var chatFormattingStripPattern = regexp.MustCompile("(?i)§[0-9A-FK-OR]")

// Name mirrors the implicit Enum.name().
func (c ChatFormatting) Name() string { return c.name }

// GetCode mirrors ChatFormatting.code.
func (c ChatFormatting) GetCode() rune { return c.code }

// ToString mirrors ChatFormatting.toString().
func (c ChatFormatting) ToString() string { return c.toString }

// String mirrors ChatFormatting.toString() (Go Stringer form).
func (c ChatFormatting) String() string { return c.toString }

// StripFormatting mirrors ChatFormatting.stripFormatting.
func StripFormatting(input *string) *string {
	if input == nil {
		return nil
	}
	stripped := chatFormattingStripPattern.ReplaceAllString(*input, "")
	return &stripped
}

// ChatFormattingGetByCode mirrors ChatFormatting.getByCode.
func ChatFormattingGetByCode(code rune) *ChatFormatting {
	sanitized := lowerRune(code)
	for _, format := range ChatFormattingValues() {
		if format.code == sanitized {
			f := format
			return &f
		}
	}
	return nil
}

func lowerRune(r rune) rune {
	if r >= 'A' && r <= 'Z' {
		return r + ('a' - 'A')
	}
	return r
}
