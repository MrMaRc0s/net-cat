package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
)

type client struct {
	conn     net.Conn
	nick     string
	room     *room
	commands chan<- command
}

func (c *client) readInput() {
	defer recoverFromPanic(c.conn.RemoteAddr())

	file, _ := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	defer file.Close()
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			log.Printf("client %s has disconnected: %v", c.conn.RemoteAddr(), err)
			c.commands <- command{id: CMD_QUIT, client: c}
			return
		}

		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/nick":
			c.commands <- command{
				id:     CMD_NICK,
				client: c,
				args:   args,
			}
		case "/join":
			c.commands <- command{
				id:     CMD_JOIN,
				client: c,
				args:   args,
			}
		case "/rooms":
			c.commands <- command{
				id:     CMD_ROOMS,
				client: c,
			}
		case "/msg":
			c.commands <- command{
				id:     CMD_MSG,
				client: c,
				args:   args,
			}
			for i := 1; i < len(args)-1; i++ {
				file.WriteString(args[i] + " ")
			}
			file.WriteString(args[len(args)-1] + "\n")
		case "/quit":
			c.commands <- command{
				id:     CMD_QUIT,
				client: c,
			}
		case "/help":
			c.commands <- command{
				id:     CMD_HELP,
				client: c,
			}
		default:
			c.commands <- command{
				id:     CMD_MSG,
				client: c,
				args:   args,
			}
			for i := 1; i < len(args)-1; i++ {
				file.WriteString(args[i] + " ")
			}
			file.WriteString(args[len(args)-1] + "\n")
		}
	}
}

func (c *client) welcome() {
	pinguin, err := os.ReadFile("pinguin.txt")
	if err != nil {
		log.Printf("error: %s", err)
	}
	entryMsg := "Welcome to TCP-Chat!\n" + string(pinguin) + "\n[ENTER THY NAME N'WAH]: "
	c.conn.Write([]byte("> " + entryMsg))

	reader := bufio.NewReader(c.conn)
	nickname, err := reader.ReadString('\n')
	if err != nil {
		c.err(err)
		return
	}
	nickname = strings.Trim(nickname, "\r\n")
	if strings.TrimSpace(nickname) == "" {
		nickname = "Anonymous"
	}

	c.commands <- command{
		id:     CMD_NICK,
		client: c,
		args:   []string{"/nick", nickname},
	}
	go c.readInput()
}

func (c *client) err(err error) {
	c.conn.Write([]byte("err: " + err.Error() + "\n"))
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}
