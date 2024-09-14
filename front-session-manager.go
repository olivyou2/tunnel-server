package main

import (
	"sync"
)

type FrontSessionManager struct {
	// sessions map[string]*FrontSession
	sessions *sync.Map
}

func newFrontSessionManager() *FrontSessionManager {
	fsm := new(FrontSessionManager)
	fsm.sessions = new(sync.Map)

	return fsm
}

func (tsm *FrontSessionManager) add(session *FrontSession) {
	// tsm.sessions[session.sessid] = session
	tsm.sessions.Store(session.sessid, session)
}

func (tsm *FrontSessionManager) get(sessId string) *FrontSession {
	val, ok := tsm.sessions.Load(sessId)

	if ok {
		return val.(*FrontSession)
	} else {
		return nil
	}

	// return tsm.sessions[sessId]

}

func (tsm *FrontSessionManager) remove(sessId string) {
	tsm.sessions.Delete(sessId)
}
