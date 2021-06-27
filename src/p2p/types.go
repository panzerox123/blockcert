package p2p

import (
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
}
