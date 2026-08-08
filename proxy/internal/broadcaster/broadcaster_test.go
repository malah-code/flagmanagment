package broadcaster

import (
	"fmt"
	"testing"
	"time"

	pb "github.com/flagmanagment/proxy/pkg/gen/sdk/v1"
	"github.com/stretchr/testify/assert"
)

func TestBroadcaster_RegisterDeregister(t *testing.T) {
	b := NewBroadcaster()

	c1 := &Client{ID: "c1", Ch: make(chan *pb.RulesetDelta, 10)}
	b.Register(c1)
	assert.Equal(t, 1, b.Count())

	b.Deregister("c1")
	assert.Equal(t, 0, b.Count())
}

func TestBroadcaster_Broadcast(t *testing.T) {
	b := NewBroadcaster()

	c1 := &Client{ID: "c1", Ch: make(chan *pb.RulesetDelta, 10)}
	c2 := &Client{ID: "c2", Ch: make(chan *pb.RulesetDelta, 10)}

	b.Register(c1)
	b.Register(c2)

	delta := &pb.RulesetDelta{Version: "v2"}
	b.Broadcast(delta)

	select {
	case msg := <-c1.Ch:
		assert.Equal(t, "v2", msg.Version)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("c1 did not receive delta")
	}

	select {
	case msg := <-c2.Ch:
		assert.Equal(t, "v2", msg.Version)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("c2 did not receive delta")
	}
}

func TestBroadcaster_SlowConsumerDrop(t *testing.T) {
	b := NewBroadcaster()

	c1 := &Client{ID: "c1", Ch: make(chan *pb.RulesetDelta, 1)}
	b.Register(c1)

	// Fill buffer
	b.Broadcast(&pb.RulesetDelta{Version: "v1"})

	// Buffer full, next broadcast should non-block drop
	done := make(chan struct{})
	go func() {
		b.Broadcast(&pb.RulesetDelta{Version: "v2"})
		close(done)
	}()

	select {
	case <-done:
		// Success — non-blocking
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Broadcast blocked on full consumer channel")
	}
}

func TestBroadcaster_500ClientsStress(t *testing.T) {
	b := NewBroadcaster()
	const numClients = 500

	clients := make([]*Client, numClients)
	for i := 0; i < numClients; i++ {
		clients[i] = &Client{
			ID: fmt.Sprintf("client-%d", i),
			Ch: make(chan *pb.RulesetDelta, 10),
		}
		b.Register(clients[i])
	}

	assert.Equal(t, numClients, b.Count())

	start := time.Now()
	b.Broadcast(&pb.RulesetDelta{Version: "stress-v1"})

	for i := 0; i < numClients; i++ {
		select {
		case msg := <-clients[i].Ch:
			assert.Equal(t, "stress-v1", msg.Version)
		case <-time.After(500 * time.Millisecond):
			t.Fatalf("client %d did not receive delta in time", i)
		}
	}

	duration := time.Since(start)
	assert.True(t, duration < 500*time.Millisecond, "Fanout took longer than 500ms")
}
