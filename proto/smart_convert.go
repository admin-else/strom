package proto

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/admin-else/strom/data"
	"github.com/admin-else/strom/event"
	"github.com/admin-else/strom/proto_base"
	"github.com/admin-else/strom/proto_generated"
)

func SmartConvertibleTo(from, to reflect.Type) bool {
	switch from.Kind() {
	case reflect.Struct:
		if to.Kind() != reflect.Struct {
			return false
		}
		if from.NumField() != to.NumField() {
			return false
		}
		for i := range from.NumField() {
			if !SmartConvertibleTo(from.Field(i).Type, to.Field(i).Type) {
				return false
			}
		}
		return true
	case reflect.Pointer, reflect.Slice:
		if from.Kind() != to.Kind() {
			return false
		}
		return SmartConvertibleTo(from.Elem(), to.Elem())
	default:
		return from.ConvertibleTo(to)
	}
}

func SmartConvert(from reflect.Value, to reflect.Type) (ret reflect.Value) {
	switch from.Kind() {
	case reflect.Struct:
		ret = reflect.New(to).Elem()
		for i := range from.NumField() {
			ret.Field(i).Set(SmartConvert(from.Field(i), to.Field(i).Type))
		}
	case reflect.Pointer:
		if from.IsNil() {
			return reflect.Zero(to)
		}
		ret = SmartConvert(from.Elem(), to.Elem()).Addr()
	case reflect.Slice:
		l := from.Len()
		ret = reflect.MakeSlice(to, l, l)
		for i := range l {
			ret.Index(i).Set(SmartConvert(from.Index(i), to.Elem()))
		}
	default:
		return from.Convert(to)
	}
	return ret
}

// RegisterUntilLatest registers a handler for all versions between the passed in packet version and the latest supported version.
// Warning: If a bad handler is passed in, the program will panic.
func (c *Conn) RegisterUntilLatest(h any) {
	c.registerUntil(h, LatestVersion, false)
}

// RegisterUntil registers a handler for all versions between the passed in packet version and the specified version.
// Warning: If a bad handler is passed in, the program will panic.
func (c *Conn) RegisterUntil(h any, until string) {
	c.registerUntil(h, until, false)
}

// RegisterCriticalUntilLatest registers a handler for all versions between the passed in packet version and the latest supported version.
// If this handler returns an error, the connection will be closed.
// Warning: If a bad handler is passed in, the program will panic.
func (c *Conn) RegisterCriticalUntilLatest(h any) {
	c.registerUntil(h, LatestVersion, true)
}

// RegisterCriticalUntil registers a handler for all versions between the passed in packet version and the specified version.
// If this handler returns an error, the connection will be closed.
// Warning: If a bad handler is passed in, the program will panic.
func (c *Conn) RegisterCriticalUntil(h any, until string) {
	c.registerUntil(h, until, true)
}

func (c *Conn) registerUntil(h any, until string, critical bool) {
	eventType, hv := event.ValidateHandler(h)
	if critical {
		c.RegisterDirect(eventType, hv)
	} else {
		c.RegisterCustomType(eventType, hv)
	}
	if !eventType.Implements(reflect.TypeFor[proto_base.EncodeDecodeAble]()) {
		panic("expected method with argument that implements proto_base.EncodeDecodeAble")
	}
	packetInfo, found := LookupPacketInfoByType(reflect.Zero(eventType).Interface().(proto_base.EncodeDecodeAble))
	if !found {
		panic("packet not found")
	}
	versionInfo, err := data.LookUpVersionByProtocolVersion(packetInfo.ProtocolVersion)
	if err != nil {
		panic("protocol version not found")
	}
	startIndex := slices.Index(proto_generated.SupportedVersions, versionInfo.MinecraftVersion)
	endIndex := slices.Index(proto_generated.SupportedVersions, until)
	for i := startIndex + 1; i <= endIndex; i++ {
		versionName := proto_generated.SupportedVersions[i]
		versionInfo, err = data.LookUpVersionByName(versionName)
		if err != nil {
			panic("protocol version not found")
		}
		var newPacketInfo proto_base.PacketInfo
		newPacketInfo, found = LookupPacketInfoByNameProtocolVersionStateAndDirection(packetInfo.Name, versionInfo.Version, packetInfo.State, packetInfo.Direction)
		if !found {
			panic("packet not found")
		}
		newT := reflect.TypeOf(newPacketInfo.Type)
		if !SmartConvertibleTo(newT, reflect.TypeOf(packetInfo.Type)) {
			panic(fmt.Sprintf("packet %v is not convertible to %v", newT, reflect.TypeOf(packetInfo.Type)))
		}
		handleFunc := func(packet any) error {
			v := hv.Call([]reflect.Value{SmartConvert(reflect.ValueOf(packet), reflect.TypeOf(packetInfo.Type))})[0]
			if v.IsNil() {
				return nil
			}
			return v.Interface().(error)
		}
		if critical {
			c.RegisterDirect(newT, reflect.ValueOf(handleFunc))
		} else {
			c.RegisterCustomType(newT, reflect.ValueOf(handleFunc))
		}
	}
}
