package goiconv

/*
#include <iconv.h>
#include <stdlib.h>
*/
import "C"
import (
	"bytes"
	"syscall"
	"unsafe"
)

type IConv struct {
	h C.iconv_t
}

func NewIConv(from, to string) (*IConv, error) {
	cfrom, cto := C.CString(from), C.CString(to)
	defer C.free(unsafe.Pointer(cfrom))
	defer C.free(unsafe.Pointer(cto))

	h, err := C.iconv_open(cto, cfrom)
	if err != nil {
		return nil, err
	}

	return &IConv{h}, nil
}

func (iconv *IConv) Close() error {
	_, err := C.iconv_close(iconv.h)
	return err
}

func (iconv *IConv) IConv(inbuf, outbuf []byte) (inleftbytes, outleftbytes []byte, err error) {
	inbufp, outbufp := (*C.char)(unsafe.Pointer(&inbuf[0])), (*C.char)(unsafe.Pointer(&outbuf[0]))
	inbufLen, outbufCap := len(inbuf), cap(outbuf)
	inbytesleft, outbytesleft := C.size_t(inbufLen), C.size_t(outbufCap)

	_, err = C.iconv(iconv.h, &inbufp, &inbytesleft, &outbufp, &outbytesleft)
	return inbuf[inbufLen-int(inbytesleft):], outbuf[:outbufCap-int(outbytesleft)], err
}

func (iconv *IConv) ConvertString(instr string) (string, error) {
	inbuf := []byte(instr)
	inbufp := (*C.char)(unsafe.Pointer(&inbuf[0]))
	inbytesleft := C.size_t(len(inbuf))
	outbufCap := len(inbuf) * 2
	if outbufCap < 10 {
		// too small outbuf makes C.iconv returning always E2BIG.
		outbufCap += 10
	}
	outbuf := make([]byte, outbufCap)
	var buf bytes.Buffer
	var err error

	for inbytesleft > 0 {
		outbufp := (*C.char)(unsafe.Pointer(&outbuf[0]))
		outbytesleft := C.size_t(cap(outbuf))

		_, err = C.iconv(iconv.h, &inbufp, &inbytesleft, &outbufp, &outbytesleft)
		buf.Write(outbuf[:outbufCap-int(outbytesleft)])
		if err != syscall.E2BIG {
			break
		}
	}

	return buf.String(), err
}
