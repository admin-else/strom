package protodef

type (
	BitFieldEntry struct {
		Name   string
		Signed bool
		Size   int
	}

	BitField []BitFieldEntry

	EntryHolderSet struct {
		BaseName  string
		Otherwise struct {
			Name string
			Type any
		}
	}

	Switch struct {
		CompareTo string
		Fields    map[string]any
		Default   any
	}

	// TypeForward is not a real type, I just didn't implement it all yet, and it's just the underlying type
	TypeForward struct {
		Type string
	}

	Buffer struct {
		Count     int
		CountType any
	}

	Array struct {
		CountType any
		Count     string
		Type      any
	}

	Container []struct {
		Name string
		Type any
		Anon bool
	}

	Mapper struct {
		Mappings map[string]string
		Type     any
	}

	RegistryEntryHolderSet struct {
		Base struct {
			Name string
			Type any
		}
		Otherwise struct {
			Name string
			Type any
		}
	}
	BitFlags struct {
		Flags []string
		Type  any
	}

	EntityMetadataLoop struct {
		EndVal int
		Type   any
	}
)
