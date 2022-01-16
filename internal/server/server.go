package server

import (
	"context"
	"log"
	"net/http"
	"poc-publisher/config"
	"poc-publisher/internal/dao"
	"poc-publisher/internal/handlers"
	"poc-publisher/internal/services"
	"strings"

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
		Set:       envConfig.DB.Set,
		NameSpace: envConfig.DB.NameSpace,
	})

	// pubsubClient := client.GetNewClient(ctx, client.PublishConfig{
	// 	ProjectID: envConfig.PubSub.ProjectID,
	// 	TopicID:   envConfig.PubSub.TopicID,
	// 	Logger:    logger,
	// })
	// defer pubsubClient.CloseClient()

	publisherService := services.GetPublisherService(&services.PublisherConfig{
		Logger: logger,
		Dao:    aeroDB,
		//PSClient: pubsubClient,
	})

	logger.Info("starting poc-publisher on...", zap.String("PORT", envConfig.DeployConfig.Port))
	router := handlers.InitRouter(ctx, handlers.RouterConfig{
		Service: publisherService,
		Logger:  logger,
	})
	certPath := strings.Join([]string{envConfig.DeployConfig.AppDir, "cert"}, "/")
	log.Fatal(http.ListenAndServeTLS(":"+envConfig.DeployConfig.Port, certPath+"/local.crt", certPath+"/local.key", router))

}
