package main

import "os"

func AssingPort() string {
	var port string
	if len(os.Args) >= 2 {
		port = ":" + os.Args[len(os.Args)-1]
	} else {
		port = ":8989"
	}
	return port
}
