package goiconv

import (
	"bytes"
	"syscall"
	"testing"
)

func iconv(from, to string, inbuf []byte, t *testing.T) (dst []byte) {
	iconv, err := NewIConv(from, to)
	if err != nil {
		t.Errorf("NewIConv(%s, %s) failed. err=%s", from, to, err)
	}
	defer iconv.Close()

	abuf := make([]byte, 0)
	outbufLen := 1
	outbuf := make([]byte, outbufLen)
	for len(inbuf) > 0 {
		inbuf, outbuf, err = iconv.IConv(inbuf, outbuf)
		if len(outbuf) > 0 {
			abuf = append(abuf, outbuf...)
		}
		if err != nil {
			if err == syscall.E2BIG {
				outbufLen++
				outbuf = make([]byte, outbufLen)
			} else {
				t.Fatalf("IConv() failed. err=%s", err)
			}
		}
	}

	return abuf
}

func TestUtf8ToUtf8(t *testing.T) {
	in := []byte("日本語")
	out := iconv("utf8", "utf8", in, t)
	if !bytes.Equal(in, out) {
		t.Error("unexpected output returned")
	}
}

func TestIConv(t *testing.T) {
	// source data. 日本語 in utf8
	inbuf := []byte("\xe6\x97\xa5\xe6\x9c\xac\xe8\xaa\x9e")

	// open iconv
	iconv, err := NewIConv("utf8", "iso-2022-jp") // from:utf8 to:iso-2022-jp
	if err != nil {
		// C.iconv_open returned an error.
		t.Fatal("unexpected error ", err)
		return
	}
	// call C.iconv_close later
	defer iconv.Close()

	// convert all bytes
	buf := make([]byte, 0)
	outbufLen := 1
	outbuf := make([]byte, outbufLen)
	for len(inbuf) > 0 {
		inbuf, outbuf, err = iconv.IConv(inbuf, outbuf)
		t.Logf("inbuf  consumed  = % x\n", inbuf)
		t.Logf("outbuf converted = % x\n", outbuf)

		if len(outbuf) > 0 {
			buf = append(buf, outbuf...)
		}
		if err != nil {
			if err == syscall.E2BIG {
				outbufLen++
				t.Logf("expanding outbuf. size = %d\n", outbufLen)
				outbuf = make([]byte, outbufLen)
			} else {
				t.Fatal("unexpected error ", err)
			}
		}
	}

	t.Logf("bytes converted to iso-2022-jp = % x\n", buf)
}
