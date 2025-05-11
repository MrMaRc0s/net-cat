package main

import (
	"net"
	"sync"
)

type room struct {
	name    string
	members map[net.Addr]*client
	sync.RWMutex
}

func (r *room) broadcast(sender *client, msg string) {

	r.RLock()
	defer r.RUnlock()
	for _, m := range r.members {
		m.msg(msg)
	}
}
