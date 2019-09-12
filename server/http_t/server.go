package http_t

import "net"

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


}

func (srv *Server)Serve(ln net.Listener) error {
	for {
		rw, _ := ln.Accept()

		// 将这个TCP请求包了一层
		srv.newConn(rw)
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
}
