package utf16conv

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"
)

func Example() {
	_, err := io.Copy(os.Stdout, New(os.Stdin))
	if err != nil {
		log.Fatalf("piping to stdout: %v", err)
	}
}

var utf8bom = []byte{0xEF, 0xBB, 0xBF}

func TestUTF16LEIsConverted(t *testing.T) {
	x := []byte{0xFF, 0xFE, 0x24, 0x00}
	b := bytes.NewBuffer(x)
	var o bytes.Buffer
	_, err := io.Copy(&o, New(b))
	if err != nil {
		t.Errorf("converting from UTF16-LE: %v\n", err)
	}
	if a := string(o.Bytes()); a != string(append(utf8bom, '\u0024')) {
		t.Errorf("expected $; actual: %x", a)
	}
}

func TestUTF16BEIsConverted(t *testing.T) {
	x := []byte{0xFE, 0xFF, 0x00, 0x24}
	b := bytes.NewBuffer(x)
	var o bytes.Buffer
	_, err := io.Copy(&o, New(b))
	if err != nil {
		t.Errorf("converting from UTF16-LE: %v\n", err)
	}
	if a := string(o.Bytes()); a != string(append(utf8bom, '\u0024')) {
		t.Errorf("expected $; actual: %x", a)
	}
}

func TestUTF8LeftAlone(t *testing.T) {
	x := []byte("$")
	b := bytes.NewBuffer(x)
	var o bytes.Buffer
	_, err := io.Copy(&o, New(b))
	if err != nil {
		t.Errorf("converting from UTF-8: %v\n", err)
	}
	if a := string(o.Bytes()); a != "$" {
		t.Errorf("expected $; actual: %x", a)
	}
}
