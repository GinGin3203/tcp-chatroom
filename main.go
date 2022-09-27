package main

import (
	"log"
	"net"

	"github.com/GinGin3203/protohackers/pkg/must"
)

func main() {

	addr := must.NotFail(net.ResolveTCPAddr("tcp", ":1025"))
	listener := must.NotFail(net.ListenTCP("tcp", addr))

	connManager := newManager()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("Connection established from: %s", conn.RemoteAddr().String())
		go connManager.Serve(conn)
	}
}
