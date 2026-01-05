package proto

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"io"
	"net"

	"github.com/admin-else/strom/crypto"
	"github.com/admin-else/strom/proto_base"
)

type RawConn struct {
	net.Conn
	R                    io.Reader
	W                    io.Writer
	CompressionThreshold int32
}

func (c *RawConn) SendRaw(rawPacketBytes []byte) (err error) {
	var packetBytes []byte
	if c.CompressionThreshold > 0 {
		packetBuffer := bytes.NewBuffer(nil)
		if int32(len(packetBytes)) >= c.CompressionThreshold {
			err = proto_base.EncodeVarInt(packetBuffer, int32(len(packetBytes)))
			if err != nil {
				return
			}
			_, err = zlib.NewWriter(packetBuffer).Write(rawPacketBytes)
			if err != nil {
				return
			}
		} else {
			err = proto_base.EncodeVarInt(packetBuffer, 0)
			if err != nil {
				return
			}
			_, err = packetBuffer.Write(rawPacketBytes)
			if err != nil {
				return
			}
		}
		packetBytes = packetBuffer.Bytes()
	} else {
		packetBytes = rawPacketBytes
	}
	err = proto_base.EncodeVarInt(c.W, int32(len(packetBytes)))
	if err != nil {
		return
	}
	_, err = c.W.Write(packetBytes)
	return
}

func (c *RawConn) ReceiveRaw() (packetBytes []byte, err error) {
	rawPacketLen, err := proto_base.DecodeVarInt(c.R)
	if err != nil {
		return
	}
	rawPacketBytes, err := io.ReadAll(io.LimitReader(c.R, int64(rawPacketLen)))
	if err != nil {
		return
	}
	if len(rawPacketBytes) != int(rawPacketLen) {
		err = errors.New("bad packet length")
		return
	}
	rawPacketBuffer := bytes.NewBuffer(rawPacketBytes)
	if c.CompressionThreshold > 0 {
		var packetLen int32
		packetLen, err = proto_base.DecodeVarInt(rawPacketBuffer)
		if err != nil {
			return
		}
		if packetLen == 0 {
			packetBytes, err = io.ReadAll(rawPacketBuffer)
			if err != nil {
				return
			}
		} else {
			var zReader io.ReadCloser
			zReader, err = zlib.NewReader(rawPacketBuffer)
			if err != nil {
				return
			}
			defer zReader.Close()
			packetBytes, err = io.ReadAll(zReader)
			if err != nil {
				return
			}
		}
	} else {
		packetBytes = rawPacketBytes
	}
	return
}

func (c *RawConn) SetSecret(sharedSecret []byte) (err error) {
	var b cipher.Block
	b, err = aes.NewCipher(sharedSecret)
	if err != nil {
		return
	}
	c.R = cipher.StreamReader{
		S: crypto.NewCFB8Decrypt(b, sharedSecret),
		R: c.Conn,
	}
	c.W = cipher.StreamWriter{
		S: crypto.NewCFB8Encrypt(b, sharedSecret),
		W: c.Conn,
	}
	return
}
