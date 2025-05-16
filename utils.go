package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"unicode"
)

func AssingPort() string {
	var port string
	if len(os.Args) == 2 {
		port = ":" + os.Args[len(os.Args)-1]
		if !isnumeric(os.Args[1]) {
			log.Printf("[USAGE]: ./TCPChat $port")
			os.Exit(1)
		}
	} else if len(os.Args) == 1 {
		port = ":8989"
	} else {
		log.Printf("[USAGE]: ./TCPChat $port")
		os.Exit(1)
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

func isnumeric(str string) bool {
	for _, char := range str {
		if !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}
