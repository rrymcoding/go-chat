package main

import (
	"net"
)

type room struct {
	name    string
	members map[net.Addr]*client
}

func (r *room) broadcast(sender *client, msg string) {
	for _, m := range r.members {
		//if sender.conn.RemoteAddr() != addr {
		m.msg(msg)
		//}
	}
}
