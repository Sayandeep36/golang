package main

import (
	"context"
	"log"
	"time"

	"Golang/common"
	"Golang/grpcclient"
	"Golang/proto"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// Connect to a server
	nc, err := nats.Connect("localhost:4222")
	if err != nil {
		log.Fatal("could not connect to nats", err.Error())
	}
	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal("could not create new jet stream", err.Error())
	}

	_, err = js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "salesforce",
		Subjects: []string{"sfdc.pubsub.in.*"},
	})
	if err != nil {
		log.Fatal("could not create NATS stream", err.Error())
	}

	if common.ReplayPreset == proto.ReplayPreset_CUSTOM && common.ReplayId == nil {
		log.Fatalf("the replayId variable must be populated when the replayPreset variable is set to CUSTOM")
	} else if common.ReplayPreset != proto.ReplayPreset_CUSTOM && common.ReplayId != nil {
		log.Fatalf("the replayId variable must not be populated when the replayPreset variable is set to EARLIEST or LATEST")
	}

	log.Printf("Creating gRPC client...")
	client, err := grpcclient.NewGRPCClient()
	if err != nil {
		log.Fatalf("could not create gRPC client: %v", err)
	}
	defer client.Close()

	log.Printf("Populating auth token...")
	err = client.Authenticate()
	if err != nil {
		client.Close()
		log.Fatalf("could not authenticate: %v", err)
	}

	log.Printf("Populating user info...")
	err = client.FetchUserInfo()
	if err != nil {
		client.Close()
		log.Fatalf("could not fetch user info: %v", err)
	}

	kv, _ := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket: "sfdc-replay-logs",
	})

	log.Printf("Making GetTopic request...")

	topics := []string{"/data/ContactChangeEvent", "/data/AccountChangeEvent"}
	for _, topicName := range topics {

		go run(topicName, client, js, kv, ctx)
	}
	select {}
}

func run(topicName string, client *grpcclient.PubSubClient, js jetstream.JetStream, kv jetstream.KeyValue, ctx context.Context) {
	var replayID, curReplayId []byte
	entry, err := kv.Get(ctx, topicName)
	if err == nil {
		replayID = entry.Value()
	} else {
		log.Printf("Error getting key %s: %v", topicName, err)
	}

	topic, err := client.GetTopic(topicName)
	if err != nil {
		client.Close()
		log.Fatalf("could not fetch topic: %v", err)
	}

	if !topic.GetCanSubscribe() {
		client.Close()
		log.Fatalf("this user is not allowed to subscribe to the following topic: %s", topic)
	}
	if replayID != nil {
		curReplayId = replayID
	} else {
		curReplayId = common.ReplayId
	}

	for {
		log.Printf("Subscribing to topic...")

		// use the user-provided ReplayPreset by default, but if the curReplayId variable has a non-nil value then assume that we want to
		// consume from a custom offset. The curReplayId will have a non-nil value if the user explicitly set the ReplayId or if a previous
		// subscription attempt successfully processed at least one event before crashing
		replayPreset := common.ReplayPreset
		if curReplayId != nil {
			replayPreset = proto.ReplayPreset_CUSTOM
		}

		// In the happy path the Subscribe method should never return, it will just process events indefinitely. In the unhappy path
		// (i.e., an error occurred) the Subscribe method will return both the most recently processed ReplayId as well as the error message.
		// The error message will be logged for the user to see and then we will attempt to re-subscribe with the ReplayId on the next iteration
		// of this for loop
		curReplayId, err = client.Subscribe(replayPreset, curReplayId, js, topicName, kv)
		if err != nil {
			log.Printf("error occurred while subscribing to topic: %v", err)
		}
	}

}
