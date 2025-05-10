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

	listener, err := net.Listen("tcp", ":8989")
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

	log.Printf("started server on :8989")

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
