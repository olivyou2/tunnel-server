package main

import "sync"

type FrontServerManager struct {
	fronts *sync.Map
}

func newFrontServerManager() *FrontServerManager {
	fsm := new(FrontServerManager)
	fsm.fronts = new(sync.Map)

	return fsm
}

func (fsm *FrontServerManager) add(alias string, fs *FrontServer) {
	fsm.fronts.Store(alias, fs)
}

func (fsm *FrontServerManager) get(alias string) *FrontServer {
	val, ok := fsm.fronts.Load(alias)

	if !ok {
		return nil
	}

	return val.(*FrontServer)
}

func (fsm *FrontServerManager) remove(alias string) {
	fsm.fronts.Delete(alias)
}
