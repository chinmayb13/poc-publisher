package server

import (
	"context"
	"log"
	"poc-publisher/internal/client"
	"poc-publisher/internal/dao"
	"poc-publisher/internal/services"

	"go.uber.org/zap"
)

func RunServer() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("logger initialization failed")
	}
	defer logger.Sync()

	ctx := context.Background()
	aeroDB := dao.GetDBConn(dao.AeroConfig{
		Host:      "localhost",
		Port:      3000,
		Logger:    logger,
		QueueSize: 64,
		LimitConn: true,
		Timeout:   1000,
	})

	pubsubClient := client.GetNewClient(ctx, client.PublishConfig{
		ProjectID: "poc-hdfc",
		TopicID:   "poc-publisher",
		Logger:    logger,
	})
	defer pubsubClient.CloseClient()

	publisherService := services.GetPublisherService(&services.PublisherConfig{
		Logger:   logger,
		Dao:      aeroDB,
		PSClient: pubsubClient,
	})

	logger.Info("starting poc--publisher...")
	publisherService.ReadWriteKeys(ctx)

}
