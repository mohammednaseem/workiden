package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
)

var (
	/*bridge   = struct {
		host *string
		port *string
	}{
		flag.String("mqtt_host", "mqtt.googleapis.com", "MQTT Bridge Host"),
		flag.String("mqtt_port", "8883", "MQTT Bridge Port"),
	}*/
	projectID = flag.String("project", "famous-palisade-356103", "GCP Project ID")
	subID     = flag.String("subName", "a-x", "Subscription name")
)

func main() {
	publish(*projectID, "", "")
	pullMsgsSync()
}

func publish(projectID, topicID, msg string) error {
	projectID = "famous-palisade-356103"
	topicID = "a-topic"
	msg = "Pub World"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub: NewClient: %v", err)
	}
	defer client.Close()

	t := client.Topic(topicID)
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		//panic(err)
		fmt.Println(err)
	}
	//fmt.Fprintf(w, "Published a message; msg ID: %v\n", id)
	fmt.Println("Published a message; msg ID: " + id)
	return nil
}

func pullMsgsSync() error {
	// projectID := "my-project-id"
	// subID := "my-sub"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, *projectID)
	if err != nil {
		//panic(err)
		fmt.Println(err)
	}
	defer client.Close()

	sub := client.Subscription(*subID)

	// Turn on synchronous mode. This makes the subscriber use the Pull RPC rather
	// than the StreamingPull RPC, which is useful for guaranteeing MaxOutstandingMessages,
	// the max number of messages the client will hold in memory at a time.
	sub.ReceiveSettings.Synchronous = true
	sub.ReceiveSettings.MaxOutstandingMessages = 100

	// Receive messages for 10 seconds, which simplifies testing.
	// Comment this out in production, since `Receive` should
	// be used as a long running operation.
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Second)
	defer cancel()

	var received int32
	err = sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		fmt.Println("Got message: " + string(msg.Data))
		received = received + 1
		msg.Ack()
	})
	if err != nil {
		//panic(err)
		fmt.Println(err)
	}
	fmt.Println(received)

	return nil
}
