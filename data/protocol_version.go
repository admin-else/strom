package data

//   {
//    "minecraftVersion": "1.21.10-rc1",
//    "version": 1073742098,
//    "dataVersion": 4555,
//    "usesNetty": true,
//    "majorVersion": "1.21",
//    "releaseType": "snapshot"
//  },

type ProtocolVersion struct {
	MinecraftVersion string
	Version          int
	DataVersion      int
	UsesNetty        bool
	MajorVersion     string
	ReleaseType      string
}

var ProtocolVersions []ProtocolVersion

func LookUpProtocolVersionByName(name string) (ret ProtocolVersion, err error) {
	for _, v := range ProtocolVersions {
		if v.MinecraftVersion == name {
			ret = v
			return
		}
	}
	err = UnknownMinecraftVersionError
	return
}

func LookUpProtocolVersionByVersion(version int) (ret ProtocolVersion, err error) {
	for _, v := range ProtocolVersions {
		if v.Version == version {
			ret = v
			return
		}
	}
	err = UnknownMinecraftVersionError
	return
}

func init() {
	must(LoadJson("minecraft-data/data/pc/common/protocolVersions.json", &ProtocolVersions))
}
