package tcp

import (
	"io"
	"kim"
	"kim/wire/endian"
)

func WriteFrame(w io.Writer, code kim.OpCode, payload []byte) error {
	if err := endian.WriteUint8(w, uint8(code)); err != nil {
		return err
	}
	if err := endian.WriteBytes(w, payload); err != nil {
		return err
	}
	return nil
}
