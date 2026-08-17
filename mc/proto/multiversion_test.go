package proto_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"git.anygate.cloud/anygatecloud/strom/mc/data"
	"git.anygate.cloud/anygatecloud/strom/mc/proto"
	"git.anygate.cloud/anygatecloud/strom/mc/proto_base"
	"git.anygate.cloud/anygatecloud/strom/mc/proto_generated"
)

func TestPacketRoundTripAllVersions(t *testing.T) {
	results := map[string][4]int{}
	for _, version := range proto_generated.SupportedVersions {
		protoVer := data.MustLookupProtocolVersion(version)
		total, encFailed, decFailed, ok := 0, 0, 0, 0
		for _, pin := range proto_generated.Packets {
			if pin.ProtocolVersion != int32(protoVer) {
				continue
			}
			total++
			b, err := proto.SimplePacketToBytes(pin.Type)
			if err != nil {
				if err == proto_base.ToDoError || strings.Contains(err.Error(), "unsupported") || strings.Contains(err.Error(), "index not found") || strings.Contains(err.Error(), "bad type") {
					encFailed++
					continue
				}
				t.Errorf("%s/%s: encode error: %v", version, pin.Name, err)
				continue
			}
			decoded, err := proto.SimpleBytesToPacket(b, int32(protoVer), pin.Direction, pin.State)
			if err != nil {
				t.Errorf("%s/%s: decode error: %v", version, pin.Name, err)
				continue
			}
			if unc, ok := decoded.(*proto.UnCodablePacket); ok {
				if unc.Err == proto_base.ToDoError {
					decFailed++
					continue
				}
				t.Errorf("%s/%s: uncodable: %v", version, pin.Name, unc.Err)
				continue
			}
			b2, err := proto.SimplePacketToBytes(decoded)
			if err != nil {
				t.Errorf("%s/%s: re-encode error: %v", version, pin.Name, err)
				continue
			}
			if len(b) != len(b2) {
				t.Errorf("%s/%s: byte length mismatch %d != %d", version, pin.Name, len(b), len(b2))
			}
			ok++
		}
		results[version] = [4]int{total, encFailed, decFailed, ok}
	}
	fmt.Println()
	for _, version := range proto_generated.SupportedVersions {
		r := results[version]
		fmt.Printf("%-8s %3d pkt  enc-skip:%-3d  dec-skip:%-3d  ok:%-3d\n", version, r[0], r[1], r[2], r[3])
	}
}

func TestSmartConvertCheck(t *testing.T) {
	panicCount := 0
	notConvertible := 0
	convertible := 0
	totalPairs := 0
	for _, fromV := range proto_generated.SupportedVersions {
		for _, toV := range proto_generated.SupportedVersions {
			if fromV == toV {
				continue
			}
			fromVer := data.MustLookupProtocolVersion(fromV)
			toVer := data.MustLookupProtocolVersion(toV)
			for _, pinFrom := range proto_generated.Packets {
				if pinFrom.ProtocolVersion != int32(fromVer) {
					continue
				}
				pinTo, found := proto.LookupPacketInfoByNameProtocolVersionStateAndDirection(
					pinFrom.Name, int32(toVer), pinFrom.State, pinFrom.Direction,
				)
				if !found {
					continue
				}
				totalPairs++
				func() {
					defer func() {
						if r := recover(); r != nil {
							panicCount++
						}
					}()
					fromT := typeInfo(pinFrom.Type)
					toT := typeInfo(pinTo.Type)
					if proto.SmartConvertibleTo(fromT, toT) {
						convertible++
					} else {
						notConvertible++
					}
				}()
			}
		}
	}
	fmt.Printf("smart convert: %d pairs, %d convertible, %d not, %d panics\n", totalPairs, convertible, notConvertible, panicCount)
}

func typeInfo(t proto_base.EncodeDecodeAble) reflect.Type {
	rt := reflect.TypeOf(t)
	for rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	return rt
}

func TestCrossVersionDecode(t *testing.T) {
	errors := 0
	panics := 0
	total := 0
	for _, toV := range proto_generated.SupportedVersions {
		toVer := data.MustLookupProtocolVersion(toV)
		for _, pin := range proto_generated.Packets {
			if pin.ProtocolVersion == int32(toVer) {
				continue
			}
			_, found := proto.LookupPacketInfoByNameProtocolVersionStateAndDirection(
				pin.Name, int32(toVer), pin.State, pin.Direction,
			)
			if !found {
				continue
			}
			b, err := proto.SimplePacketToBytes(pin.Type)
			if err != nil {
				continue
			}
			total++
			func() {
				defer func() {
					if r := recover(); r != nil {
						panics++
					}
				}()
				decoded, err := proto.SimpleBytesToPacket(b, int32(toVer), pin.Direction, pin.State)
				if err != nil {
					errors++
					return
				}
				if _, ok := decoded.(*proto.UnCodablePacket); ok {
					return
				}
				_, err = proto.SimplePacketToBytes(decoded)
				if err != nil {
					errors++
				}
			}()
		}
	}
	fmt.Printf("cross-version decode: %d attempts, %d errors, %d panics\n", total, errors, panics)
}
