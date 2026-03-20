// NOTICE

// Project Name: Cloaq
// Copyright © 2026 Neil Talap and/or its designated Affiliates.

// This software is licensed under the Dragonfly Public License (DPL) 1.0.

// All rights reserved. The names "Neil Talap" and any associated logos or branding
// are trademarks of the Licensor and may not be used without express written permission,
// except as provided in Section 7 of the License.

// For commercial licensing inquiries or permissions beyond the scope of this
// license, please create an issue in github.

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
