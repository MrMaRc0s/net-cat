package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func AssingPort() string {
	var port string
	if len(os.Args) == 2 {
		port = ":" + os.Args[len(os.Args)-1]
	} else {
		port = ":8989"
	}
	return port
}

func DeleteFile(filename string) {
	err := os.Remove(filename)
	if err != nil {
		fmt.Println("Error deleting file:", err)
	} else {
		fmt.Println("File deleted successfully!")
	}
}

func recoverFromPanic(clientAddr net.Addr) {
	if r := recover(); r != nil {
		log.Printf("client %s disconnected abruptly: %v", clientAddr, r)
	}
}

func handleShutdown(stop chan os.Signal, listener net.Listener, s *server) {
	<-stop
	log.Println("shutting down...")
	listener.Close()
	close(s.commands)
	s.wg.Wait()
	os.Exit(0)
}
