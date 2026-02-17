package drivers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub"
)

type PubSubDriver struct {
	projectID      string
	topicName      string
	publishEnabled bool
	client         *pubsub.Client
	topic          *pubsub.Topic
}

func NewPubSubDriver(projectID, topicName string) (*PubSubDriver, error) {
	if projectID == "" {
		projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
		if projectID == "" {
			projectID = os.Getenv("GCP_PROJECT")
		}
	}
	if projectID == "" {
		return nil, fmt.Errorf("GOOGLE_CLOUD_PROJECT no está configurado")
	}

	if topicName == "" {
		topicName = os.Getenv("PUBSUB_TOPIC_NAME")
		if topicName == "" {
			topicName = "digital-edd-sdk"
		}
	}

	publishEnabled := os.Getenv("SDKTRACKING_PUBLISH") != "false"

	return &PubSubDriver{
		projectID:      projectID,
		topicName:      topicName,
		publishEnabled: publishEnabled,
	}, nil
}

func (d *PubSubDriver) ensureClient() error {
	if d.client != nil {
		return nil
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, d.projectID)
	if err != nil {
		return err
	}

	d.client = client
	d.topic = client.Topic(d.topicName)
	fmt.Printf("[digital-edd-logger] PubSub conectado al topic: %s\n", d.topicName)
	return nil
}

func (d *PubSubDriver) Send(record map[string]interface{}) (string, error) {
	if !d.publishEnabled {
		return "publish-disabled", nil
	}

	if err := d.ensureClient(); err != nil {
		return "", err
	}

	data, err := json.Marshal(record)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	result := d.topic.Publish(ctx, &pubsub.Message{
		Data: data,
	})

	id, err := result.Get(ctx)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (d *PubSubDriver) Close() error {
	if d.client != nil {
		return d.client.Close()
	}
	return nil
}
