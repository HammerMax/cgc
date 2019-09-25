package http_t

import (
	"bufio"
	"context"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

type ResponseWriter interface {
	Header() Header

	Write([]byte) (int, error)

	WriteHeader(statusCode int)
}

// http请求的处理者
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}

// Server 定义了运行一个HTTP服务的参数
type Server struct {
	Addr    string  // TCP address to listen on, 监听的端口
	Handler Handler // handler to invoke, http.DefaultServeMux if nil. 请求的处理者，nil默认为DefaultServeMux

	ReadTimeout time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout time.Duration


}

type serverHandler struct {
	srv *Server
}

func (sh serverHandler) ServeHTTP(rw ResponseWriter, req *Request) {
	handler := sh.srv.Handler
	if handler == nil {
		handler = DefaultServeMux
	}
	// URI为"*", Method为"OPTIONS"的请求，走不到我们定义的ServeMux
	if req.RequestURI == "*" && req.Method == "OPTIONS" {
		handler = globalOptionsHandler{}
	}
	handler.ServeHTTP(rw, req)
}

func ListenAndServe(addr string, handler Handler) error {
	server := &Server{Addr: addr, Handler: handler}
	return server.ListenAndServe()
}

func (srv *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}
	return srv.Serve(ln)
}

func (srv *Server)Serve(ln net.Listener) error {
	for {
		rw, _ := ln.Accept()

		// 将这个TCP请求包了一层
		c := srv.newConn(rw)
		go c.serve()
	}
}

func (srv *Server) newConn(rwc net.Conn) *conn {
	c := &conn{
		srv: srv,
		rwc: rwc,
	}
	return c
}

// A conn represents the server side of an HTTP connection.
type conn struct {
	srv *Server
	rwc net.Conn

	bufr *bufio.Reader
	bufw *bufio.Writer
}

func (c *conn) serve() {
	defer func() {
		if err := recover(); err != nil {

		}
	}()



}

func (c *conn) readRequest(ctx context.Context) (w *response, err error) {
	req, err := readRequest(c.bufr, keepHostHeader)
	if err != nil {
	}

	w = &response{
		conn: c,
		req: req,
		reqBody: req.Body,
	}
	return
}

// 承担了解析URL并找到最合适Handler，将w, r交给Handler处理
type ServeMux struct {
	mu	sync.RWMutex
	m map[string]muxEntry
	// slice of entries sorted from longest to shortest.
	// 猜测这个排序是为了更快找到匹配的URL。因为匹配URL是按照根到子查询，所以不能用map匹配
	es []muxEntry
	hosts bool
}

type muxEntry struct {
	h Handler
	pattern string
}

// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request) {
	h, _ := mux.Handler(r)
	h.ServeHTTP(w, r)
}

// 根据req，返回对应的Handler
func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string) {
	return mux.handler("", r.URL.Path)
}

func (mux *ServeMux) handler(host, path string) (h Handler, pattern string) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	h, pattern = mux.match(path)
	return
}

// Find a handler on a handler map given a path string.
// Most-specific (longest) pattern wins
func (mux *ServeMux) match(path string) (h Handler, pattern string) {
	// Check for exact match first
	v, ok := mux.m[path]
	if ok {
		return v.h, v.pattern
	}

	// Check for longest valid match.  mux.es contains all patterns
	// that end in / sorted from longest to shortest
	// 这里就是体现为什么ServeMux es[]muxEntry 是按从长到短排序的。
	// 最匹配URL的pattern，会最先被选中
	for _, e := range mux.es {
		if strings.HasPrefix(path, e.pattern) {
			return e.h, e.pattern
		}
	}
	return nil, ""
}

// NewServeMux allocates and returns a new ServeMux.
func NewServeMux() *ServeMux { return new(ServeMux) }

// DefaultServeMux is the default ServeMux used by Serve.
var DefaultServeMux = &defaultServeMux

var defaultServeMux ServeMux

type HandlerFunc func(ResponseWriter, *Request)

// A response represents the server side of an HTTP response
type response struct {
	conn *conn
	req *Request
	reqBody io.ReadCloser

	w *bufio.Writer
}

func (w *response) Write(data []byte) (n int, err error) {
	return w.write(len(data), data, "")
}

func (w *response) write(lenData int, dataB []byte, dataS string) (n int, err error) {
	return w.w.Write(dataB)
}

// globalOptionsHandler responds to "OPTIONS *" requests.
type globalOptionsHandler struct{}

func (globalOptionsHandler) ServeHTTP(w ResponseWriter, r *Request) {
	w.Header().Set("Content-Length", "0")
}