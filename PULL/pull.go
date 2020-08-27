package main

import (
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/pubsub"
)

func pullMsgs() error {
	pID := "cloud-test-287516"
	sID := "test-cloud"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, pID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	// Consume 10 messages.
	var mu sync.Mutex
	received := 0
	sub := client.Subscription(sID)
	cctx, cancel := context.WithCancel(ctx)
	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()
		fmt.Println("Got message: " + string(msg.Data))
		msg.Ack()
		received++
		if received == 10 {
			cancel()
		}
	})
	if err != nil {
		return fmt.Errorf("Receive: %v", err)
	}
	return nil
}

func main() {
	pullMsgs()
}
