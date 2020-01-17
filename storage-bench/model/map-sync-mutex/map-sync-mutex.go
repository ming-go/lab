package mapsyncmutex

import "sync"

type mapSyncMutex struct {
	db  map[string][]byte
	mux sync.Mutex
}

func (msm *mapSyncMutex) Create(k string, v string) {

}

func (msm *mapSyncMutex) Delete(k string) {

}

func (msm *mapSyncMutex) Update(k string, v []byte) []byte {
	return []byte("")
}

func (msm *mapSyncMutex) Retrieve(k string) {

}
