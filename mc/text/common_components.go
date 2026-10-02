package text

// CommonComponents mirrors net.minecraft.network.chat.CommonComponents.
var (
	CommonComponentsEMPTY                      Component = ComponentEmpty()
	CommonComponentsOPTION_ON                  Component = ComponentTranslatable("options.on")
	CommonComponentsOPTION_OFF                 Component = ComponentTranslatable("options.off")
	CommonComponentsGUI_DONE                   Component = ComponentTranslatable("gui.done")
	CommonComponentsGUI_CANCEL                 Component = ComponentTranslatable("gui.cancel")
	CommonComponentsGUI_YES                    Component = ComponentTranslatable("gui.yes")
	CommonComponentsGUI_REMOVE                 Component = ComponentTranslatable("gui.remove")
	CommonComponentsGUI_NO                     Component = ComponentTranslatable("gui.no")
	CommonComponentsGUI_OK                     Component = ComponentTranslatable("gui.ok")
	CommonComponentsGUI_PROCEED                Component = ComponentTranslatable("gui.proceed")
	CommonComponentsGUI_CONTINUE               Component = ComponentTranslatable("gui.continue")
	CommonComponentsGUI_BACK                   Component = ComponentTranslatable("gui.back")
	CommonComponentsGUI_TO_TITLE               Component = ComponentTranslatable("gui.toTitle")
	CommonComponentsGUI_ACKNOWLEDGE            Component = ComponentTranslatable("gui.acknowledge")
	CommonComponentsGUI_OPEN_IN_BROWSER        Component = ComponentTranslatable("chat.link.open")
	CommonComponentsGUI_COPY_TO_CLIPBOARD      Component = ComponentTranslatable("chat.copy")
	CommonComponentsGUI_COPY_LINK_TO_CLIPBOARD Component = ComponentTranslatable("gui.copy_link_to_clipboard")
	CommonComponentsGUI_DISCONNECT             Component = ComponentTranslatable("menu.disconnect")
	CommonComponentsGUI_RETURN_TO_MENU         Component = ComponentTranslatable("menu.returnToMenu")
	CommonComponentsTRANSFER_CONNECT_FAILED    Component = ComponentTranslatable("connect.failed.transfer")
	CommonComponentsCONNECT_FAILED             Component = ComponentTranslatable("connect.failed")
	CommonComponentsNEW_LINE                   Component = ComponentLiteral("\n")
	CommonComponentsNARRATION_SEPARATOR        Component = ComponentLiteral(". ")
	CommonComponentsELLIPSIS                   Component = ComponentLiteral("...")
	CommonComponentsSPACE                      Component = CommonComponentsSpace()
)

// CommonComponentsSpace mirrors CommonComponents.space().
func CommonComponentsSpace() *MutableComponent { return ComponentLiteral(" ") }

// CommonComponentsDays mirrors CommonComponents.days.
func CommonComponentsDays(value int64) *MutableComponent {
	return ComponentTranslatableArgs("gui.days", value)
}

// CommonComponentsHours mirrors CommonComponents.hours.
func CommonComponentsHours(value int64) *MutableComponent {
	return ComponentTranslatableArgs("gui.hours", value)
}

// CommonComponentsMinutes mirrors CommonComponents.minutes.
func CommonComponentsMinutes(value int64) *MutableComponent {
	return ComponentTranslatableArgs("gui.minutes", value)
}

// CommonComponentsOptionStatus mirrors CommonComponents.optionStatus(boolean).
func CommonComponentsOptionStatus(value bool) Component {
	if value {
		return CommonComponentsOPTION_ON
	}
	return CommonComponentsOPTION_OFF
}

// CommonComponentsDisconnectButtonLabel mirrors CommonComponents.disconnectButtonLabel.
func CommonComponentsDisconnectButtonLabel(isLocalServer bool) Component {
	if isLocalServer {
		return CommonComponentsGUI_RETURN_TO_MENU
	}
	return CommonComponentsGUI_DISCONNECT
}

// CommonComponentsOptionStatusName mirrors CommonComponents.optionStatus(Component, boolean).
func CommonComponentsOptionStatusName(name Component, value bool) *MutableComponent {
	key := "options.off.composed"
	if value {
		key = "options.on.composed"
	}
	return ComponentTranslatableArgs(key, name)
}

// CommonComponentsOptionNameValue mirrors CommonComponents.optionNameValue.
func CommonComponentsOptionNameValue(name, value Component) *MutableComponent {
	return ComponentTranslatableArgs("options.generic_value", name, value)
}

// CommonComponentsJoinForNarration mirrors CommonComponents.joinForNarration.
func CommonComponentsJoinForNarration(components ...Component) *MutableComponent {
	result := ComponentEmpty()
	for i, component := range components {
		result.AppendComponent(component)
		if i != len(components)-1 {
			result.AppendComponent(CommonComponentsNARRATION_SEPARATOR)
		}
	}
	return result
}

// CommonComponentsJoinLines mirrors CommonComponents.joinLines(Component...).
func CommonComponentsJoinLines(lines ...Component) Component {
	return CommonComponentsJoinLinesList(lines)
}

// CommonComponentsJoinLinesList mirrors CommonComponents.joinLines(Collection).
func CommonComponentsJoinLinesList(lines []Component) Component {
	return ComponentUtilsFormatListComponents(lines, CommonComponentsNEW_LINE)
}
