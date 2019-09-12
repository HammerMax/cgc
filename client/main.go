package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("开始阅读")
	var read []byte
	conn.Read(read)
	fmt.Println(string(read))

	fmt.Println("写入完成")

}
