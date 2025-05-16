package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

func AssingPort() string {
	var port string
	if len(os.Args) == 2 {
		port_num, err := strconv.Atoi(os.Args[1])
		if err != nil || port_num < 0 || port_num > 65535 {
			log.Printf("port must be a number and between 0 and 65535")
			os.Exit(1)
		}
		port = ":" + os.Args[len(os.Args)-1]
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
