package http_t

import (
	"bufio"
	"errors"
	"golang.org/x/net/http/httpguts"
	"io"
	"strconv"
	"sync"
)

// 根据HTTP版本以及Header判断，发送或接收request、response后，conn是否要关闭
func shouldClose(major, minor int, header Header, removeCloseHeader bool) bool {
	// HTTP版本小于1，应该没有keepalive的功能
	if major < 1 {
		return true
	}

	conv := header["Connection"]
	hasClose := httpguts.HeaderValuesContainsToken(conv, "close")
	if major == 1 && minor == 0 {
		return hasClose || !httpguts.HeaderValuesContainsToken(conv, "keep-alive")
	}

	if hasClose && removeCloseHeader {
		header.Del("Connection")
	}

	return hasClose
}

// 解析Request的中间Reader
type transferReader struct {
	// Input
	Header Header
	StatusCode int
	RequestMethod string
	ProtoMajor int
	ProtoMinor int

	// Output
	Body io.ReadCloser
	ContentLength int64
	TransferEncoding []string
	Close bool
	Trailer	Header
}

// msg is *Request or *Response
func readTransfer(msg interface{}, r *bufio.Reader) (err error) {
	t := &transferReader{RequestMethod: "GET"}

	switch rr := msg.(type) {
	case *Request:
		 t.Header = rr.Header
		 t.RequestMethod = rr.Method
		 t.ProtoMajor = rr.ProtoMajor
		 t.ProtoMinor = rr.ProtoMinor
		 t.StatusCode = 200
		 t.Close = rr.Close
	default:
		panic("unexpected type")
	}

	// Default to HTTP/1.1
	if t.ProtoMajor == 0 && t.ProtoMinor == 0 {
		t.ProtoMajor, t.ProtoMinor = 1, 1
	}

	realLength, err := strconv.Atoi(t.Header["Content-Length"][0])

	t.Body = &body{src: io.LimitReader(r, int64(realLength)), closing: t.Close}

	switch rr := msg.(type) {
	case *Request:
		rr.Body = t.Body
		rr.ContentLength = t.ContentLength
	}

	return nil
}

// body turns a Reader into a ReadCloser.
// Close ensures that the body has been fully read
// and then reads the trailer if necessary.
type body struct {
	src          io.Reader
	hdr          interface{}   // non-nil (Response or Request) value means read trailer
	r            *bufio.Reader // underlying wire-format reader for the trailer
	closing      bool          // is the connection to be closed after reading body?
	doEarlyClose bool          // whether Close should stop early

	mu         sync.Mutex // guards following, and calls to Read and Close
	sawEOF     bool
	closed     bool
	earlyClose bool   // Close called and we didn't read to the end of src
	onHitEOF   func() // if non-nil, func to call when EOF is Read
}

var ErrBodyReadAfterClose = errors.New("http: invalid Read on closed Body")

func (b *body) Read(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return 0, ErrBodyReadAfterClose
	}
	return b.readLocked(p)
}

func (b *body) readLocked(p []byte) (n int, err error) {
	if b.sawEOF {
		return 0, io.EOF
	}
	n, err = b.src.Read(p)

	if err == io.EOF {
		b.sawEOF = true

		if b.hdr != nil {

		} else {

		}

	}

	return
}

func (b *body) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return nil
	}
	var err error
	b.closed = true
	return err
}