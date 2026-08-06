package packet_info

import (
	"flag"
	"fmt"
	"strconv"

	"git.anygate.cloud/anygatecloud/strom/mc/data"
	"git.anygate.cloud/anygatecloud/strom/mc/proto"
	"git.anygate.cloud/anygatecloud/strom/mc/proto_base"
	"git.anygate.cloud/anygatecloud/strom/mc/proto_generated"
)

var cmd = flag.NewFlagSet("packet-info", flag.ContinueOnError)
var versionFlag = cmd.String("version", "", "protocol version or minecraft version name")
var stateFlag = cmd.String("state", "", "handshaking, status, login, configuration, play")
var directionFlag = cmd.String("direction", "", "toserver, toclient")
var idFlag = cmd.Int("id", -1, "packet id")
var nameFlag = cmd.String("name", "", "packet name")

func Run(args []string) (err error) {
	err = cmd.Parse(args)
	if err != nil {
		return
	}

	var version int32
	if *versionFlag != "" {
		v, err := data.LookUpVersionByName(*versionFlag)
		if err != nil {
			v2, err := parseVersion(*versionFlag)
			if err != nil {
				return fmt.Errorf("invalid version: %s", *versionFlag)
			}
			v, err = data.LookUpVersionByProtocolVersion(v2)
			if err != nil {
				return err
			}
		}
		version = v.Version
	}

	var state proto_base.State
	if *stateFlag != "" {
		state, err = parseState(*stateFlag)
		if err != nil {
			return
		}
	}

	var direction proto_base.Direction
	if *directionFlag != "" {
		direction, err = parseDirection(*directionFlag)
		if err != nil {
			return
		}
	}

	if *idFlag != -1 && *nameFlag != "" {
		err = fmt.Errorf("cannot specify both id and name")
		return
	}

	if *idFlag != -1 {
		if *versionFlag == "" || *stateFlag == "" || *directionFlag == "" {
			err = fmt.Errorf("version, state, and direction are required when looking up by id")
			return
		}
		p, ok := proto.LookUpTypeByPacketInfo(direction, state, int32(*idFlag), version)
		if !ok {
			err = fmt.Errorf("packet not found")
			return
		}
		printPacketInfo(p)
		return
	}

	if *nameFlag != "" {
		if *versionFlag == "" {
			err = fmt.Errorf("version is required when looking up by name")
			return
		}
		if state != 0 || direction != 0 {
			p, ok := proto.LookupPacketInfoByNameProtocolVersionStateAndDirection(*nameFlag, version, state, direction)
			if !ok {
				err = fmt.Errorf("packet not found")
				return
			}
			printPacketInfo(p)
		} else {
			var found bool
			for _, p := range proto_generated.Packets {
				if p.Name == *nameFlag && p.ProtocolVersion == version {
					printPacketInfo(p)
					found = true
				}
			}
			if !found {
				err = fmt.Errorf("packet not found")
			}
		}
		return
	}

	err = fmt.Errorf("must specify either -id or -name")
	return
}

func parseState(s string) (proto_base.State, error) {
	switch s {
	case "handshaking":
		return proto_base.Handshaking, nil
	case "status":
		return proto_base.Status, nil
	case "login":
		return proto_base.Login, nil
	case "configuration":
		return proto_base.Configuration, nil
	case "play":
		return proto_base.Play, nil
	default:
		return 0, fmt.Errorf("unknown state: %s", s)
	}
}

func parseDirection(s string) (proto_base.Direction, error) {
	switch s {
	case "toserver":
		return proto_base.ToServer, nil
	case "toclient":
		return proto_base.ToClient, nil
	default:
		return 0, fmt.Errorf("unknown direction: %s", s)
	}
}

func parseVersion(val string) (int32, error) {
	v, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(v), nil
}

func printPacketInfo(p proto_base.PacketInfo) {
	v, _ := data.LookUpVersionByProtocolVersion(p.ProtocolVersion)
	fmt.Printf("Name:       %s\n", p.Name)
	fmt.Printf("ID:         %d\n", p.PacketId)
	fmt.Printf("State:      %s\n", stateToString(p.State))
	fmt.Printf("Direction:  %s\n", p.Direction)
	fmt.Printf("Version:    %d (%s)\n", p.ProtocolVersion, v.MinecraftVersion)
}

func stateToString(s proto_base.State) string {
	switch s {
	case proto_base.Handshaking:
		return "handshaking"
	case proto_base.Status:
		return "status"
	case proto_base.Login:
		return "login"
	case proto_base.Configuration:
		return "configuration"
	case proto_base.Play:
		return "play"
	default:
		return "unknown"
	}
}
