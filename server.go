package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type server struct {
	rooms    map[string]*room
	commands chan command
	wg       sync.WaitGroup
	sync.RWMutex
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client)
		case CMD_HELP:
			s.help(cmd.client)
		}
	}
}

func (s *server) newClient(conn net.Conn) *client {
	log.Printf("new client has connected: %s", conn.RemoteAddr().String())

	s.wg.Add(1)

	return &client{
		conn:     conn,
		commands: s.commands,
	}
}

func (s *server) nick(c *client, args []string) {
	if len(args) != 2 || args[1] == "" {
		c.msg("nick is required. usage: /nick NAME")
		return
	}
	if len(args[1]) > 20 {
		c.msg("nick is too large, the limit is 20 characters")
		return
	}

	c.nick = args[1]
	c.msg(fmt.Sprintf("you will be called %s", c.nick))
}

func (s *server) join(c *client, args []string) {
	if len(args) < 2 {
		c.msg("room name is required. usage: /join ROOM_NAME")
		return
	}
	if len(args[1]) > 30 {
		c.msg("room name is too large, the limit is 30 characters")
		return
	}

	roomName := args[1]

	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s has joined the room", c.nick))

	c.msg(fmt.Sprintf("welcome to %s", r.name))
}

func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("available are: %s", strings.Join(rooms, ", ")))
}

func (s *server) msg(c *client, args []string) {
	file, _ := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer file.Close()

	msg := strings.Join(args[1:], " ")
	msgline := "[" + time.Now().Format("2006-01-02 15:04:05") + "]" + "[" + c.nick + "]" + ": " + msg
	if c.room != nil {
		c.room.broadcast(c, msgline)
		file.WriteString(msgline + "\n")
	} else {
		c.msg("join a room first n'wah")
	}
}

func (s *server) quit(c *client) {
	log.Printf("client has left the chat: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.msg("sad to see you go :(")
	c.conn.Close()

	//s.wg.Done()
}

func (s *server) help(c *client) {
	c.msg("the list of all available commands are:\n /nick\n /rooms\n /join\n /quit")
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}
