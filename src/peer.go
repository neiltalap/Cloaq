// © 2026 Polypodomy — AGPLv3

package network

import (
	"crypto/ecdh"
	"encoding/hex"
	"net"
	"sync"
)

type Peer struct {
	PublicKey *ecdh.PublicKey
	Addr      *net.Addr
}

type PeerTable struct {
	mu    sync.RWMutex
	peers map[string]*Peer
}

func NewPeerTable() *PeerTable {
	return &PeerTable{
		peers: make(map[string]*Peer),
	}
}

func pubKeyHex(pub *ecdh.PublicKey) string {
	if pub == nil {
		return ""
	}
	return hex.EncodeToString(pub.Bytes())
}

func (pt *PeerTable) AddPeer(p *Peer) {
	if pt == nil || p == nil || p.PublicKey == nil || p.Addr == nil {
		return
	}

	key := pubKeyHex(p.PublicKey)

	pt.mu.Lock()
	defer pt.mu.Unlock()

	if pt.peers == nil {
		pt.peers = make(map[string]*Peer)
	}

	pt.peers[key] = p
}

func (pt *PeerTable) GetPeer(pubKey *ecdh.PublicKey) *Peer {
	if pt == nil || pubKey == nil {
		return nil
	}

	key := pubKeyHex(pubKey)

	pt.mu.RLock()
	defer pt.mu.RUnlock()

	return pt.peers[key]
}
