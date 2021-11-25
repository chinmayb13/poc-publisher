package client

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
)

//go:generate mockery --name=PubSubClient
type PubSubClient interface {
	PublishMessage(ctx context.Context, msg string) error
	CloseClient()
}

type PublishConfig struct {
	ProjectID string
	TopicID   string
	Logger    *zap.Logger
}

type psClient struct {
	*pubsub.Client
	logger *zap.Logger
	topic  *pubsub.Topic
}

func GetNewClient(ctx context.Context, config PublishConfig) PubSubClient {
	client, err := pubsub.NewClient(ctx, config.ProjectID)
	if err != nil {
		log.Fatalf("client creation failed %s", err.Error())
	}
	topic := client.Topic(config.TopicID)
	return &psClient{
		Client: client,
		logger: config.Logger,
		topic:  topic,
	}
}

func (client *psClient) CloseClient() {
	client.Close()
}

func (client *psClient) PublishMessage(ctx context.Context, msg string) error {
	client.logger.Info("publishing to pubsub", zap.String("msg", msg))
	result := client.topic.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),		
	})
	id, err := result.Get(ctx)
	if err != nil {
		client.logger.Error("failed publishing the message", zap.String("err", err.Error()))
		return err
	}
	client.logger.Info("Published the msg", zap.String("msg", msg), zap.String("id", id))
	return nil
}
