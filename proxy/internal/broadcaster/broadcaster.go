package broadcaster

import (
	"sync"

	pb "github.com/flagmanagment/proxy/pkg/gen/sdk/v1"
)

type Client struct {
	ID string
	Ch chan *pb.RulesetDelta
}

type Broadcaster struct {
	mu      sync.RWMutex
	clients map[string]*Client
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		clients: make(map[string]*Client),
	}
}

func (b *Broadcaster) Register(c *Client) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.clients[c.ID] = c
}

func (b *Broadcaster) Deregister(clientID string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if c, ok := b.clients[clientID]; ok {
		delete(b.clients, clientID)
		close(c.Ch)
	}
}

func (b *Broadcaster) Broadcast(delta *pb.RulesetDelta) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, c := range b.clients {
		select {
		case c.Ch <- delta:
		default:
			// Non-blocking drop if channel buffer is full to prevent slow consumer stalling
		}
	}
}

func (b *Broadcaster) Count() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.clients)
}
