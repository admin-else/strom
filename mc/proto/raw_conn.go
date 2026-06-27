package proto

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"io"
	"net"
	"sync"

	"github.com/admin-else/strom/mc/crypto"
	"github.com/admin-else/strom/mc/proto_base"
)

type RawConn struct {
	net.Conn
	R                       io.Reader
	W                       io.Writer
	compressionThreshold    int32
	CompressionThresholdMut sync.RWMutex
}

func (c *RawConn) SetCompressionThreshold(threshold int32) {
	c.CompressionThresholdMut.Lock()
	defer c.CompressionThresholdMut.Unlock()
	c.compressionThreshold = threshold
}

func (c *RawConn) GetCompressionThreshold() int32 {
	c.CompressionThresholdMut.RLock()
	defer c.CompressionThresholdMut.RUnlock()
	return c.compressionThreshold
}

func (c *RawConn) SendRaw(rawPacketBytes []byte) (err error) {
	//fmt.Println("send:", rawPacketBytes)
	var packetBytes []byte
	threshold := c.GetCompressionThreshold()
	if threshold >= 0 {
		packetBuffer := bytes.NewBuffer(nil)
		if int32(len(rawPacketBytes)) >= threshold {
			err = proto_base.EncodeVarInt(packetBuffer, int32(len(rawPacketBytes)))
			if err != nil {
				return
			}
			zW := zlib.NewWriter(packetBuffer) // ignore we don't care about the result if it errors
			_, err = zW.Write(rawPacketBytes)
			if err != nil {
				return
			}
			err = zW.Close()
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
	packetWithLenBuffer := bytes.NewBuffer(nil)
	err = proto_base.EncodeVarInt(packetWithLenBuffer, int32(len(packetBytes)))
	if err != nil {
		return
	}
	_, err = packetWithLenBuffer.Write(packetBytes)
	if err != nil {
		return
	}
	_, err = c.W.Write(packetWithLenBuffer.Bytes())
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
	if c.GetCompressionThreshold() >= 0 {
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
			packetBytes, err = io.ReadAll(zReader)
			if err != nil {
				return
			}
			err = zReader.Close()
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
