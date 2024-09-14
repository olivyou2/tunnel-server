package main

import (
	"net"
)

type FrontServer struct {
	listener *net.Listener
	sm       *SessManager

	alias       string
	fsm         *FrontSessionManager
	managerSess *TunnelSession
}

func createFrontServer(sm *SessManager, host string, managerSess *TunnelSession, alias string) *FrontServer {
	listener, err := net.Listen("tcp", host)

	if nil != err {
		logging("Failed to listen front server")
		return nil
	}

	frontServer := new(FrontServer)

	frontServer.managerSess = managerSess
	frontServer.listener = &listener
	frontServer.alias = alias
	frontServer.fsm = new(FrontSessionManager)
	frontServer.sm = sm

	return frontServer
}

func (frontServer *FrontServer) accept() {
	listener := *frontServer.listener

	for {
		conn, err := listener.Accept()

		if nil != err {
			logging("Failed to accept connection")
			return
		}

		go frontServer.frontSocketProcessing(&conn)
	}
}

func (frontServer *FrontServer) frontSocketProcessing(conn *net.Conn) {
	// Request to create new connection

	sessId := createUuid()

	bw := newBufferWriter()
	bw.writeString("newConnection")
	bw.writeString(sessId)

	pw := newBufferWriter()
	pw.writeFixBuffer(bw.getBytes())

	frontSession := new(FrontSession)

	frontSession.conn = conn
	frontSession.sessid = sessId
	frontSession.alias = frontServer.alias

	frontServer.fsm.add(frontSession)

	manager := frontServer.managerSess
	managerConn := *manager.conn

	_, err := managerConn.Write(pw.getBytes())

	if nil != err {
		logging("[ERROR] handshake write failed [", manager.sessid, manager.agentType, "]")
		logging("[ERROR]", err)
	}

	// start frontRecvProcessing when tunnel handshake successed.
}

func (frontServer *FrontServer) frontRecvProcessing(sessId string) {
	buffer := make([]byte, 1024)
	// conn := frontServer.sm.frontSessManager.get(sessId).conn
	conn := frontServer.fsm.get(sessId).conn
	tunnelSess := frontServer.sm.tunnelSessManager.get(sessId)

	for {
		n, err := (*conn).Read(buffer)

		logging("[DATA] FRONT -> TUNNEL [", sessId, n, "bytes readed]")

		if nil != err {
			logging("[ERROR] Front Session Closed", sessId, tunnelSess.agentType)
			(*tunnelSess.conn).Close()
			frontServer.sm.tunnelSessManager.remove(sessId)
			return
		}

		if n == 0 {
			(*tunnelSess.conn).Close()
			return
		}

		_, werr := (*tunnelSess.conn).Write(buffer[:n])

		if nil != werr {
			logging(werr)
			return
		}
	}
}
