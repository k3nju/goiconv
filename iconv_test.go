package goiconv

import (
	"syscall"
	"testing"
)

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
