package main

import (
	"io"
	"log"
	"net"
)

func OnConnect(conn net.Conn) (err error) {
	target, err := net.Dial("tcp", "127.0.0.1:25566")
	if err != nil {
		return
	}
	errChan := make(chan error)
	go func() {
		_, innerErr := io.Copy(target, conn)
		errChan <- innerErr
	}()
	go func() {
		_, innerErr := io.Copy(conn, target)
		errChan <- innerErr
	}()
	return <-errChan
}

func OnConnectWrapper(conn net.Conn) {
	defer conn.Close()
	err := OnConnect(conn)
	if err != nil {
		log.Println("Error: ", err)
	}
}

func startListener() (err error) {
	var l net.Listener
	l, err = net.Listen("tcp", ":25565")
	if err != nil {
		return
	}
	defer l.Close()
	for {
		var client net.Conn
		client, err = l.Accept()
		if err != nil {
			break
		}
		go OnConnectWrapper(client)
	}
	return
}

func main() {
	err := startListener()
	if err != nil {
		panic(err)
	}
}
