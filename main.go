package main

import (
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {
	s := newServer()
	go s.run()

	port := AssingPort()

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}
	defer listener.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	go func() {
		<-stop
		log.Println("shutting down...")
		listener.Close()
		close(s.commands)
		s.wg.Wait()
		os.Exit(0)
	}()

	log.Printf("started server on %s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("unable to accept connection: %s", err.Error())
			return
		}

		c := s.newClient(conn)

		go c.welcome()
	}
}
