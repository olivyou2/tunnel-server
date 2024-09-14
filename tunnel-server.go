package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
)

type TunnelServer struct {
	listener *net.Listener
	sm       *SessManager
}

func logging(a ...any) {
	if false {
		fmt.Println(a...)
	}
}

func createTunnelServer(sm *SessManager, address string) *TunnelServer {
	listener, err := net.Listen("tcp", address)

	if nil != err {
		logging("Failed to create tunnel server")
		return nil
	}

	tunnelServer := new(TunnelServer)
	tunnelServer.listener = &listener
	tunnelServer.sm = sm

	return tunnelServer
}

func (tunnelServer *TunnelServer) accept() {
	for {
		conn, err := (*tunnelServer.listener).Accept()

		tunnelSession := new(TunnelSession)
		tunnelSession.agentType = unset
		tunnelSession.conn = &conn

		if nil != err {
			logging("Failed to accept connection")
			return
		}

		go tunnelServer.tunnelRecvProcessing(tunnelSession)
	}
}

func (tunnelServer *TunnelServer) tunnelPacketProcessing(sess *TunnelSession, buffer []byte) {
	br := newBufferReader(buffer)

	if sess.agentType == unset {
		agentType := br.readString()
		agentId := br.readString()

		sess.sessid = agentId
		tunnelServer.sm.tunnelSessManager.add(sess)

		if agentType == "manager" {
			sess.agentType = manager
			tunnelServer.sm.tunnelSessManager.manager = sess

			alias := br.readString()
			sess.tunnelAlias = alias

			fs := createFrontServer(tunnelServer.sm, ":0", sess, alias)
			fs.fsm = newFrontSessionManager()

			listener := *fs.listener
			port := listener.Addr().(*net.TCPAddr).Port

			fmt.Println(alias, "Tunnel open at", port)
			sess.sendListenOk("localhost:" + strconv.Itoa(port))

			go fs.accept()

			tunnelServer.sm.frontServerManager.add(alias, fs)

			logging("Manager SESS-ID", agentId)
		} else if agentType == "agent" {
			sess.agentType = agent

			alias := br.readString()

			fs := tunnelServer.sm.frontServerManager.get(alias)

			frontSession := fs.fsm.get(agentId)
			sess.frontSession = frontSession

			logging("AGENT POLLED : [", sess.sessid, "]")

			go fs.frontRecvProcessing(agentId)
		}

	}
}

func (tunnelServer *TunnelServer) tunnelRecvProcessing(sess *TunnelSession) {

	buffer := make([]byte, 1024)
	bw := newBufferWriter()
	conn := *sess.conn
	readTimes := 0

	for {
		readTimes += 1
		n, err := conn.Read(buffer)

		if nil != err {
			logging("Recv Error [", sess.agentType, ",", sess.sessid, "]", err)

			// When tunnel disconnected, then front socket must be closed
			if nil != sess.frontSession {
				(*sess.frontSession.conn).Close()
				fs := tunnelServer.sm.frontServerManager.get(sess.frontSession.alias)
				fs.fsm.remove(sess.sessid)

				// tunnelServer.sm.frontSessManager.remove(sess.sessid)
			}

			if sess.agentType == manager {
				tunnelServer.sm.tunnelSessManager.remove(sess.sessid)

				frontListener := *tunnelServer.sm.frontServerManager.get(sess.tunnelAlias).listener
				frontListener.Close()

				tunnelServer.sm.frontServerManager.remove(sess.tunnelAlias)
			}

			return
		}

		if n > 0 {
			if sess.agentType != unset {
				logging("[DATA] ARRIVAL FROM TUNNEL [", sess.agentType, ",", sess.sessid, "]")
			}
			if sess.agentType == agent {
				// Tunnel to front
				(*sess.frontSession.conn).Write(buffer[:n])

				logging("[DATA] TUNNEL -> FRONT [", sess.sessid, "]")
				continue
			}

			bw.writeBuffer(buffer[:n])

			for {
				if len(bw.getBytes()) < 4 {
					break
				}

				br := newBufferReader(bw.getBytes())
				packetSize := br.readInt()

				if packetSize+4 <= int32(br.reader.Size()) {
					packet := make([]byte, packetSize)

					br.reader.Read(packet)
					tunnelServer.tunnelPacketProcessing(sess, packet)

					cropped, err := io.ReadAll(br.reader)
					if err != nil {
						logging(err)
					}

					bw = newBufferWriter()
					bw.writeBuffer(cropped)
				} else {
					break
				}
			}
		} else if n == 0 {
			break
		}
	}
}
