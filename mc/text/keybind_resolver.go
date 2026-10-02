package text

// KeybindResolver mirrors net.minecraft.network.chat.contents.KeybindResolver.
// The Java static Function<String, Supplier<Component>> is modelled as a
// package-global hook.
var keybindKeyResolver = func(name string) func() Component {
	return func() Component { return ComponentLiteral(name) }
}

// KeybindResolverSetKeyResolver mirrors KeybindResolver.setKeyResolver.
func KeybindResolverSetKeyResolver(resolver func(name string) func() Component) {
	keybindKeyResolver = resolver
}
