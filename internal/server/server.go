package server

import (
	"context"
	"log"
	"poc-publisher/config"
	"poc-publisher/internal/client"
	"poc-publisher/internal/dao"
	"poc-publisher/internal/services"
	"poc-publisher/internal/utils"

	"go.uber.org/zap"
)

func RunServer() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("logger initialization failed")
	}
	defer logger.Sync()

	envConfig, err := config.LoadConfig("../..")
	if err != nil {
		logger.Fatal("error reading env file", zap.String("err", err.Error()))
	}

	ctx := context.Background()
	aeroDB := dao.GetDBConn(dao.AeroConfig{
		Host:      envConfig.DB.Host,
		Port:      envConfig.DB.Port,
		Logger:    logger,
		QueueSize: envConfig.DB.QueueSize,
		LimitConn: envConfig.DB.LimitConn,
		Timeout:   envConfig.DB.Timeout,
	})

	pubsubClient := client.GetNewClient(ctx, client.PublishConfig{
		ProjectID: envConfig.PubSub.ProjectID,
		TopicID:   envConfig.PubSub.TopicID,
		Logger:    logger,
	})
	defer pubsubClient.CloseClient()

	publisherService := services.GetPublisherService(&services.PublisherConfig{
		Logger:   logger,
		Dao:      aeroDB,
		PSClient: pubsubClient,
	})

	logger.Info("starting poc--publisher...")
	publisherService.ReadWriteKeys(ctx, utils.GetRandomSequence(1, 1000)-1)

}
