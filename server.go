package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}
func (s *server) start() {

	listener, err := net.Listen("tcp", ":12000")
	if err != nil {
		log.Fatalf("[notice]\t unable to start server: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("[notice]\t server started on %s", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[notice]\t failed to accept connection: %s", err.Error())
			continue
		}

		c := s.newClient(conn)
		s.send_welcome(c)
		var default_nick []string
		default_nick = append(default_nick, "/nick")
		default_nick = append(default_nick, "<bot"+conn.RemoteAddr().String()+">")
		s.nick(c, default_nick)
		var default_room []string
		default_room = append(default_room, "/join")
		default_room = append(default_room, "#general")
		s.join(c, default_room)

		go c.readInput()
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
		case CMD_ELEVATE:
			s.elevate(cmd.client, cmd.args)
		case CMD_EXEC:
			s.exec(cmd.client, cmd.args)
		case CMD_SEND_ALL:
			s.broadcast_to_all_rooms(cmd.client, cmd.args)
		case CMD_LIST:
			s.list_room_members(cmd.client)
		case CMD_MENU:
			s.display_menu(cmd.client)
		}
	}
}

func (s *server) newClient(conn net.Conn) *client {
	log.Printf("new client has joined: %s", conn.RemoteAddr().String())

	return &client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
		isSuper:  false,
	}
}

func (s *server) nick(c *client, args []string) {
	if len(args) < 2 {
		c.msg("nick is required. usage: /nick NAME")
		return
	}

	c.nick = args[1]
	c.srv_msg(fmt.Sprintf("all right, I will call you %s", c.nick))
}

func (s *server) join(c *client, args []string) {
	if len(args) < 2 {
		c.msg("room name is required. usage: /join ROOM_NAME")
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

	r.broadcast(c, fmt.Sprintf("%s joined the room", c.nick))

	c.srv_msg(fmt.Sprintf("welcome to <%s>", roomName))
}

func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("available rooms: %s", strings.Join(rooms, ", ")))
}
func (s *server) getRoomsList() {

	println("Rooms:")

	for _, r := range s.rooms {

		println(r.name)
	}
}

func (s *server) msg(c *client, args []string) {
	if len(args) < 2 {
		c.msg("message is required, usage: /msg MSG")
		return
	}

	msg := strings.Join(args[1:], " ")
	c.room.broadcast(c, c.nick+": "+msg)
}
func (s *server) elevate(c *client, args []string) {
	if len(args) < 2 {
		c.msg("password is required, usage: /elevate PASSWORD")
		return
	}
	password := args[1]

	if password == "oou812" {
		c.isSuper = true
		c.msg("** Welcome Super **")
	} else {
		c.msg("password is incorrect, please use the correct password")
	}

	//msg := strings.Join(args[1:], " ")
	//c.room.broadcast(c, c.nick+": "+msg)
}
func (s *server) exec(c *client, args []string) {

	if c.isSuper {

		// get the target client
		if len(args) < 2 {
			c.msg("a target nick is required, usage: /exec TARGET")
			return
		}
		target := args[1]

		t := s.rooms[c.room.name]

		for _, m := range t.members {

			if m.nick == target {

				// then execute the command on that client
				m.exec_client_command(c)
				break
			}
		}

	}

}

func (s *server) quit(c *client) {
	log.Printf("client has left the chat: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.msg("sad to see you go =(")
	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}
func (s *server) list_room_members(c *client) {

	if len(c.room.members) > 1 {

		c.srv_msg("members [" + strconv.Itoa(len(c.room.members)) + "] in [" + c.room.name + "]")
		for _, m := range c.room.members {

			if m.nick != c.nick {
				c.srv_msg(m.nick)
			}

		}
	} else {

		c.srv_msg("you are the only one in this room")
	}

}
func (s *server) broadcast_to_all_rooms(c *client, args []string) {

	if c.isSuper {
		if len(args) < 2 {
			c.msg("message is required, usage: /msg MSG")
			return
		}

		msg := strings.Join(args[1:], " ")
		for _, value := range s.rooms {

			arr := value

			for _, m := range arr.members {

				arr.broadcast(m, "<*super* "+c.nick+" *super*>"+msg)
			}
		}
	}
}
func (s *server) send_welcome(c *client) {

	c.conn.Write([]byte(msg_welcome + "\n\n"))
	s.display_menu(c)

}
func (s *server) display_menu(c *client) {

	c.conn.Write([]byte("--System Menu--\n"))
	c.conn.Write([]byte("commands\n"))
	c.conn.Write([]byte("/join\n"))
	c.conn.Write([]byte("/nick\n"))
	c.conn.Write([]byte("/rooms\n"))
	c.conn.Write([]byte("/list\n"))
	c.conn.Write([]byte("/quit\n"))
	c.conn.Write([]byte("--------------\n"))
	c.conn.Write([]byte("\n"))
	c.conn.Write([]byte("\n"))

}
