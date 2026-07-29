package data

//   {
//    "minecraftVersion": "1.21.10-rc1",
//    "version": 1073742098,
//    "dataVersion": 4555,
//    "usesNetty": true,
//    "majorVersion": "1.21",
//    "releaseType": "snapshot"
//  },

type VersionInfo struct {
	MinecraftVersion string
	Version          int32
	DataVersion      int32
	UsesNetty        bool
	MajorVersion     string
	ReleaseType      string
}

var ProtocolVersions []VersionInfo

// LookUpVersionByName returns the VersionInfo for the given Minecraft version name (e.g. "1.21.11").
func LookUpVersionByName(name string) (ret VersionInfo, err error) {
	for _, v := range ProtocolVersions {
		if v.MinecraftVersion == name {
			ret = v
			return
		}
	}
	err = UnknownMinecraftVersionError
	return
}

// LookUpVersionByProtocolVersion returns the VersionInfo for the given protocol version number.
func LookUpVersionByProtocolVersion(version int32) (ret VersionInfo, err error) {
	for _, v := range ProtocolVersions {
		if v.Version == version {
			ret = v
			return
		}
	}
	err = UnknownMinecraftVersionError
	return
}

// LookUpVersionByDataVersion returns the VersionInfo for the given Minecraft data version number.
func LookUpVersionByDataVersion(version int32) (ret VersionInfo, err error) {
	for _, v := range ProtocolVersions {
		if v.DataVersion == version {
			ret = v
			return
		}
	}
	err = UnknownMinecraftVersionError
	return
}

// MustLookupProtocolVersion panics if the version is not found only intended for testing code.
func MustLookupProtocolVersion(version string) int32 {
	v, err := LookUpVersionByName(version)
	must(err)
	return v.Version
}

func init() {
	must(LoadJson("minecraft-data/pc/common/protocolVersions.json", &ProtocolVersions))
}
