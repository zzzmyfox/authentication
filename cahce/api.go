package cahce

import (
	"authentication/cahce/consistenthash"
	"sync"
)

const defaultReplicas = 50

type Api struct {
	self     string
	basePath string
	mu       sync.Mutex // guards peers and httpGetters
	peers    *consistenthash.Map
}

func (a *Api) Set(peers ...string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.peers = consistenthash.New(defaultReplicas, nil)
	a.peers.Add(peers...)
}
