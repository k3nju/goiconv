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

func TestUtf8ToISO2022JP(t *testing.T) {
	in := []byte("日本語")
	out := iconv("utf8", "iso-2022-jp", in, t)
	iso2022jp := []byte("\x1b\x24\x42\x46\x7c\x4b\x5c\x38\x6c")
	if !bytes.Equal(out, iso2022jp) {
		t.Error("unexpected output returned")
	}
}

func TestIConv(t *testing.T) {
	// source data. 日本語 in utf8
	inbuf := []byte("\xe6\x97\xa5\xe6\x9c\xac\xe8\xaa\x9e")

	// open iconv
	iconv, err := NewIConv("utf8", "sjis") // from:utf8 to:sjis
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

	t.Logf("bytes converted to sjis = % x\n", buf)
	sjis := []byte("\x93\xfa\x96\x7b\x8c\xea")
	if !bytes.Equal(buf, sjis) {
		t.Error("converting from utf8 to sjis failed")
	}
}

func TestConvertBytes(t *testing.T) {
	// 日本語 in utf8
	inbuf := []byte("\xe6\x97\xa5\xe6\x9c\xac\xe8\xaa\x9e")

	// open iconv
	iconv, err := NewIConv("utf8", "euc-jp")
	if err != nil {
		t.Fatal("unexpected error ", err)
		return
	}
	defer iconv.Close()

	buf, err := iconv.ConvertBytes(inbuf)
	if err != nil {
		t.Fatal("unexpected error ", err)
		return
	}

	eucjp := []byte("\xc6\xfc\xcb\xdc\xb8\xec")
	if !bytes.Equal(buf, eucjp) {
		t.Error("converting from utf8 from euc-jp failed")
	}
}
