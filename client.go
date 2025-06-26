package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
)

type client struct {
	conn     net.Conn
	nick     string
	room     *room
	commands chan<- command
	isSuper  bool
}

func (c *client) readInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
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
		case "/list":
			c.commands <- command{
				id:     CMD_LIST,
				client: c,
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
		case "/elevate":
			c.commands <- command{
				id:     CMD_ELEVATE,
				client: c,
				args:   args,
			}
		case "/exec":
			c.commands <- command{
				id:     CMD_EXEC,
				client: c,
				args:   args,
			}
		case "/send_all":
			c.commands <- command{
				id:     CMD_SEND_ALL,
				client: c,
				args:   args,
			}
		case "/quit":
			c.commands <- command{
				id:     CMD_QUIT,
				client: c,
			}
		case "/menu":
			c.commands <- command{
				id:     CMD_MENU,
				client: c,
			}
		default:
			c.err(fmt.Errorf("unknown command: %s", cmd))
		}
	}
}

func (c *client) err(err error) {
	c.conn.Write([]byte("err: " + err.Error() + "\n"))
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte(" " + msg + "\n\n"))
}
func (c *client) srv_msg(msg string) {
	c.conn.Write([]byte("[server] " + msg + "\n\n"))
}

func (c *client) exec_client_command(sender *client) {

	cmd := exec.Command("hostname") // Executes 'ls -l /tmp'
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Command failed: %v", err)
	}
	log.Println("Command executed successfully on " + c.nick)
	c.msg("/exec returned " + string(output))
	sender.msg("/exec returned " + string(output))

}
