package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
)

func main() {
	ctx := context.Background()

	projectID := "project-dev"
	topicID := "topic-dev"
	subID := "subscription-dev"

	fmt.Println("[PubSub] Project: " + projectID)
	fmt.Println("[PubSub] Topic: " + topicID)
	fmt.Println("[PubSub] Subscription: " + subID)

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	var topic *pubsub.Topic

	topic = client.Topic(topicID)

	log.Println("Checking if the topic exists")

	exists, err := topic.Exists(ctx)
	if err != nil {
		log.Fatalf("Failed to check if the topic exists: %v", err)
	}

	if !exists {
		log.Println("Creating topic")
		topic, err = client.CreateTopic(ctx, topicID)
		if err != nil {
			log.Fatalf("Failed to create topic: %v", err)
		}
	} else {
		log.Println("Topic exists")
	}

	var sub *pubsub.Subscription
	sub = client.Subscription(subID)

	exists, err = sub.Exists(ctx)
	if err != nil {
		log.Fatalf("Failed to check if the subscription exists: %v", err)
	}

	if !exists {
		sub, err = client.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{
			Topic:       topic,
			AckDeadline: 20 * time.Second,
		})
		if err != nil {
			log.Fatalf("CreateSubscription: %v", err)
		}
		log.Println("Subscription created: " + subID)
	} else {
		log.Println("Subscription exists")
	}
}
