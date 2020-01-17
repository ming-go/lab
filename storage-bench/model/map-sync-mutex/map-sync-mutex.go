package mapsyncmutex

import "sync"

type mapSyncMutex struct {
	db  map[string][]byte
	mux sync.RWMutex
}

func (msm *mapSyncMutex) Create(k string, v []byte) {
	msm.mux.Lock()
	defer msm.mux.Unlock()
	msm.db[k] = v
}

func (msm *mapSyncMutex) Delete(k string) {
	msm.mux.Lock()
	defer msm.mux.Unlock()
	delete(msm.db, k)
}

func (msm *mapSyncMutex) Update(k string, v []byte) {
	msm.mux.Lock()
	defer msm.mux.Unlock()
	msm.db[k] = v
}

func (msm *mapSyncMutex) Retrieve(k string) []byte {
	msm.mux.RLock()
	defer msm.mux.RUnlock()
	return msm.db[k]
}
