package p2p

import (
	"sync"

	"github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/panzerox123/blockcert/src/certificate"
)

type P2pNode struct {
	node                   host.Host
	blockchain             *certificate.BlockChain
	pubsub                 *pubsub.PubSub
	blockchainTopic        *pubsub.Topic
	blockchainSubscription *pubsub.Subscription
	newcertTopic           *pubsub.Topic
	newcertSubscription    *pubsub.Subscription
	interfaces             []string
	Status                 P2pStatus
	LockNet                sync.Mutex
}

type P2pStatus struct {
	Status     string
	LockStatus sync.Mutex
}

type NewCertPublish struct {
	Data       []byte
	PrivateKey string
}

var DISABLE_DISCOVERY bool = false

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"
