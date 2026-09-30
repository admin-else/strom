package raw_to_tmcpr

import (
	"bytes"
	"compress/zlib"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/admin-else/strom/mc/proto/replay"
	"github.com/admin-else/strom/mc/proto_base"
)

var Cmd = flag.NewFlagSet("raw-to-tmcpr", flag.ContinueOnError)

var ClientboundIn = Cmd.String("clientbound", "", "raw tcpflow file for clientbound (server->client) direction")
var ServerboundIn = Cmd.String("serverbound", "", "raw tcpflow file for serverbound (client->server) direction")
var ClientboundOut = Cmd.String("clientbound-out", "clientbound.tmcpr", "output TMCPR file for clientbound")
var ServerboundOut = Cmd.String("serverbound-out", "serverbound.tmcpr", "output TMCPR file for serverbound")
var ClientboundState = Cmd.String("clientbound-state", "login", "initial state for clientbound: handshake, login, config, status, play")
var ServerboundState = Cmd.String("serverbound-state", "handshake", "initial state for serverbound: handshake, login, config, status, play")
var UseConfig = Cmd.Bool("config", false, "use Configuration state transitions (for protocol >= 764, Minecraft 1.20.2+)")

type captureReader struct {
	r              io.Reader
	out            *replay.Writer
	state          proto_base.State
	threshold      int32
	startTime      time.Time
	haveStartTime  bool
	direction      proto_base.Direction
	pendingThreshold int32
	useConfig      bool
}

func (cr *captureReader) setStartTime() {
	if !cr.haveStartTime {
		cr.startTime = time.Now()
		cr.haveStartTime = true
	}
}

func (cr *captureReader) timestamp() uint32 {
	return uint32(time.Since(cr.startTime).Milliseconds())
}

func readPacketBytes(r io.Reader, threshold int32) ([]byte, error) {
	packetLen, err := proto_base.DecodeVarInt(r)
	if err != nil {
		return nil, fmt.Errorf("read packet length: %w", err)
	}

	rawBytes := make([]byte, packetLen)
	_, err = io.ReadFull(r, rawBytes)
	if err != nil {
		return nil, fmt.Errorf("read packet data: %w", err)
	}

	if threshold >= 0 {
		buf := bytes.NewReader(rawBytes)
		uncompressedLen, err := proto_base.DecodeVarInt(buf)
		if err != nil {
			return nil, fmt.Errorf("read uncompressed length: %w", err)
		}
		if uncompressedLen == 0 {
			return io.ReadAll(buf)
		} else {
			zReader, err := zlib.NewReader(buf)
			if err != nil {
				return nil, fmt.Errorf("zlib reader: %w", err)
			}
			defer zReader.Close()
			decompressed, err := io.ReadAll(zReader)
			if err != nil {
				return nil, fmt.Errorf("decompress: %w", err)
			}
			return decompressed, nil
		}
	}

	return rawBytes, nil
}

func (cr *captureReader) trackState(packetBytes []byte) {
	buf := bytes.NewReader(packetBytes)
	packetId, err := proto_base.DecodeVarInt(buf)
	if err != nil {
		return
	}

	switch cr.state {
	case proto_base.Handshaking:
		if cr.direction == proto_base.ToServer && packetId == 0 {
			_, err = proto_base.DecodeVarInt(buf) // protocol version
			if err != nil {
				return
			}
			_, err = proto_base.DecodeString(buf) // server address
			if err != nil {
				return
			}
			portBytes := make([]byte, 2)
			_, err = io.ReadFull(buf, portBytes)
			if err != nil {
				return
			}
			nextState, err := proto_base.DecodeVarInt(buf)
			if err != nil {
				return
			}
			if nextState == 2 {
				cr.state = proto_base.Login
			} else if nextState == 1 {
				cr.state = proto_base.Status
			}
		}
	case proto_base.Login:
		if cr.direction == proto_base.ToClient {
			if packetId == 0x03 {
				threshold, err := proto_base.DecodeVarInt(buf)
				if err == nil {
					cr.threshold = threshold
				}
			} else if packetId == 0x02 {
				if cr.useConfig {
					cr.state = proto_base.Configuration
				} else {
					cr.state = proto_base.Play
				}
			}
		} else if cr.direction == proto_base.ToServer && packetId == 0 {
			if cr.pendingThreshold >= 0 {
				cr.threshold = cr.pendingThreshold
			}
		}
	case proto_base.Configuration:
		if cr.direction == proto_base.ToClient && packetId == 0x02 {
			cr.state = proto_base.Play
		}
	case proto_base.Status:
	case proto_base.Play:
	}
}

func (cr *captureReader) process() (count int, err error) {
	for {
		packetBytes, err := readPacketBytes(cr.r, cr.threshold)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				return count, nil
			}
			return count, err
		}

		cr.setStartTime()
		cr.trackState(packetBytes)

		err = cr.out.WriteRaw(packetBytes, cr.timestamp())
		if err != nil {
			return count, fmt.Errorf("write TMCPR: %w", err)
		}
		count++
	}
}

func Run(args []string) (err error) {
	err = Cmd.Parse(args)
	if err != nil {
		return
	}

	if *ClientboundIn == "" && *ServerboundIn == "" {
		return fmt.Errorf("at least one of -clientbound or -serverbound must be specified")
	}

	parseState := func(name string) (proto_base.State, error) {
		switch name {
		case "handshake":
			return proto_base.Handshaking, nil
		case "login":
			return proto_base.Login, nil
		case "config":
			return proto_base.Configuration, nil
		case "play":
			return proto_base.Play, nil
		case "status":
			return proto_base.Status, nil
		default:
			return 0, fmt.Errorf("unknown state: %s (valid: handshake, login, config, status, play)", name)
		}
	}

	var compressionThreshold int32 = -1

	if *ClientboundIn != "" {
		state, err := parseState(*ClientboundState)
		if err != nil {
			return err
		}
		compressionThreshold, err = processFile(*ClientboundIn, *ClientboundOut, state, proto_base.ToClient, -1, *UseConfig)
		if err != nil {
			return fmt.Errorf("clientbound: %w", err)
		}
	}

	if *ServerboundIn != "" {
		state, err := parseState(*ServerboundState)
		if err != nil {
			return err
		}
		_, err = processFile(*ServerboundIn, *ServerboundOut, state, proto_base.ToServer, compressionThreshold, *UseConfig)
		if err != nil {
			return fmt.Errorf("serverbound: %w", err)
		}
	}

	return
}

func processFile(inPath, outPath string, initial proto_base.State, dir proto_base.Direction, initialThreshold int32, useConfig bool) (finalThreshold int32, err error) {
	in, err := os.Open(inPath)
	if err != nil {
		return initialThreshold, fmt.Errorf("open input: %w", err)
	}
	defer in.Close()

	out, err := os.Create(outPath)
	if err != nil {
		return initialThreshold, fmt.Errorf("create output: %w", err)
	}
	defer out.Close()

	cr := &captureReader{
		r:         in,
		out:       replay.NewWriter(out),
		state:     initial,
		threshold: -1,
		direction: dir,
		pendingThreshold: initialThreshold,
		useConfig: useConfig,
	}
	if dir == proto_base.ToClient {
		cr.threshold = initialThreshold
	}

	count, err := cr.process()
	if err != nil {
		return cr.threshold, err
	}

	fmt.Printf("%s: %d packets written to %s (threshold=%d)\n", inPath, count, outPath, cr.threshold)
	return cr.threshold, nil
}
