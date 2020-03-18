package main

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/net/http/httpguts"
	"log"
	"net"
	"net/textproto"
	"strings"
)

func indexFunc(r rune) bool {
	fmt.Println(string(r))
	return false
}

type s struct {
	str string
}

func main() {
	//b := strings.HasSuffix("hello sss/sd", "sd")
	//fmt.Println(b)

	ln, err := net.Listen("tcp", ":9191")
	if err != nil {
		log.Panic(err)
	}

	conn, err := ln.Accept()
	if err != nil {
		log.Panic(err)
	}

	b := make([]byte, 1000)
	n, err := conn.Read(b)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("read n:%d, body:%s", n, b)

	buf := bytes.NewBuffer(b)
	bufReader := bufio.NewReader(buf)
	textReader := textproto.NewReader(bufReader)
	line, err := textReader.ReadLine()
	line1, err := textReader.ReadMIMEHeader()

	fmt.Printf("line:%s\n", line)
	fmt.Printf("line1:%s", line1)
}

func isNotToken(r rune) bool {
	return !httpguts.IsTokenRune(r)
}

func validMethod(method string) bool {
	/*
	     Method         = "OPTIONS"                ; Section 9.2
	                    | "GET"                    ; Section 9.3
	                    | "HEAD"                   ; Section 9.4
	                    | "POST"                   ; Section 9.5
	                    | "PUT"                    ; Section 9.6
	                    | "DELETE"                 ; Section 9.7
	                    | "TRACE"                  ; Section 9.8
	                    | "CONNECT"                ; Section 9.9
	                    | extension-method
	   extension-method = token
	     token          = 1*<any CHAR except CTLs or separators>
	*/
	return len(method) > 0 && strings.IndexFunc(method, isNotToken) == -1
}
