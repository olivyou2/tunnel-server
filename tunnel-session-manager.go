package main

import "sync"

type TunnelSessionManager struct {
	sessions *sync.Map
	manager  *TunnelSession
}

func newTunnelSessionManager() *TunnelSessionManager {
	tsm := new(TunnelSessionManager)
	tsm.sessions = new(sync.Map)

	return tsm
}

func (tsm *TunnelSessionManager) add(session *TunnelSession) {
	// tsm.sessions[session.sessid] = session
	tsm.sessions.Store(session.sessid, session)
}

func (tsm *TunnelSessionManager) get(sessId string) *TunnelSession {
	// return tsm.sessions[sessId]
	val, ok := tsm.sessions.Load(sessId)

	if ok {
		return val.(*TunnelSession)
	} else {
		return nil
	}
}

func (tsm *TunnelSessionManager) remove(sessId string) {
	// delete(tsm.sessions, sessId)
	tsm.sessions.Delete(sessId)
}
