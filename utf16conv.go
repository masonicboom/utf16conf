// Package utf16conv automatically converts UTF-16 to UTF-8.
package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"unicode/utf16"
	"unicode/utf8"
)

// utf16to8 translates a UTF-16 byte array to a UTF-8 io.Reader.
func utf16to8(u16s []uint16) bytes.Buffer {
	runebuf := make([]byte, 4)
	var buf bytes.Buffer
	rus := utf16.Decode(u16s)
	for _, ru := range rus {
		nr := utf8.EncodeRune(runebuf, ru)
		_, err := buf.Write(runebuf[:nr])
		if err != nil {
			log.Printf("converting utf-16le -> utf-8: %v\n", err)
		}
	}
	return buf
}

// buf16To8 translates a UTF-16 buffer of specified endian-ness to a UTF-8 io.Reader.
func buf16to8(buf bytes.Buffer, bo binary.ByteOrder) io.Reader {
	bs := buf.Bytes()
	u16s := make([]uint16, len(bs)/2)
	err := binary.Read(&buf, bo, u16s)
	if err != nil {
		log.Printf("uh oh: %v", err)
		return nil
	}
	buf = utf16to8(u16s)
	return &buf
}

// New wraps a UTF-16 input source and translates it to UTF-8.
// Returns an empty io.Reader if r is unreadable.
// Returns an empty io.Reader if r is not UTF-16.
func New(r io.Reader) io.Reader {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	if err != nil {
		log.Printf("wrapping reader with utf16conv: %v", err)
		return new(bytes.Buffer)
	}

	bs := buf.Bytes()
	if bs[0] == 0xFF && bs[1] == 0xFE {
		// UTF-16LE (little-endian)
		r = buf16to8(buf, binary.LittleEndian)
		if r == nil {
			return new(bytes.Buffer)
		}
		return r
	} else if bs[0] == 0xFE && bs[1] == 0xFF {
		// UTF-16BE (big-endian)
		r = buf16to8(buf, binary.BigEndian)
		if r == nil {
			return new(bytes.Buffer)
		}
		return r
	}

	return &buf
}
