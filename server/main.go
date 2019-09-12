package main

import (
	"bytes"
	"net/http"
	"time"
)

func main() {
	http.ListenAndServe()
	b := bytes.Buffer{}
	NewBuffer
	ticker:=time.NewTicker(60 * time.Second)

	for {
		<-ticker.C


		go hanle()
	}



}
