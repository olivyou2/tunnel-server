package main

import "net"

type FrontSession struct {
	conn   *net.Conn
	sessid string
	alias  string
}
