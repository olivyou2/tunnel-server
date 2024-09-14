package main

import "net"

type AgentType string

const (
	unset   AgentType = "unset"
	manager AgentType = "manager"
	agent   AgentType = "agent"
)

type TunnelSession struct {
	conn      *net.Conn
	agentType AgentType
	sessid    string

	frontSession *FrontSession
	tunnelAlias  string
}

func (sess *TunnelSession) sendListenOk(host string) {
	bw := newBufferWriter()
	bw.writeString("listenOk")
	bw.writeString(host)

	pw := newBufferWriter()
	pw.writeFixBuffer(bw.getBytes())

	(*sess.conn).Write(pw.getBytes())
}
