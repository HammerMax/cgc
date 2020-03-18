package http_t

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"testing"
)

func TestReadRequest(t *testing.T) {
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Panic(err)
	}

	conn, err := ln.Accept()
	if err != nil {
		log.Panic(err)
	}

	bufReader := bufio.NewReader(conn)
	req, err := readRequest(bufReader, false)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(req)
	fmt.Println(req.URL.Query().Get("name"))
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Printf("body:%s", body)
}
