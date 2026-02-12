package node

import (
	"log"
	//"cloaq/src/network"
)

type NodeID struct {
	rootKey     string   // long-term (anonymization) storage key
	sessionKeys []string // per session key
	circuitKeys []string // one-time one way key
}

func Bootstrap() bool {
	log.Println("Bootstrapping")
	// Default peer table that we will be using to bootstrap our network
	// for loop to try and connect to each one one at a time
	// if error return false

	return true
}

func PeerDiscovery() {

}

func CreateListener() {

}
