package main

import (
	"log"
	"net"
	"os"
	"os/signal"
)

var filename string = "chatHistory"

func main() {
	s := newServer()
	go s.run()

	port := AssingPort()

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Printf("Failed to start server on %s: server have opened on default port", port)
		listener, err = net.Listen("tcp", ":8989")
		port = ":8989"
		if err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}
	defer listener.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go handleShutdown(stop, listener, s)

	defer DeleteFile(filename)

	log.Printf("started server on %s", port)

	for {
		//fmt.Println(runtime.NumGoroutine(), "active goroutines")
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("unable to accept connection: %s", err.Error())
			return
		}

		c := s.newClient(conn)

		go c.welcome()
	}
}
