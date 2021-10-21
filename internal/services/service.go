package services

import (
	"context"
	"poc-publisher/internal/client"
	"poc-publisher/internal/dao"
	"poc-publisher/internal/utils"

	"go.uber.org/zap"
)

type PublisherService interface {
	ReadWriteKeys(ctx context.Context) error
}

type publisher struct {
	logger   *zap.Logger
	dao      dao.AeroDBService
	psClient client.PubSubClient
}

type PublisherConfig struct {
	Logger   *zap.Logger
	Dao      dao.AeroDBService
	PSClient client.PubSubClient
}

func GetPublisherService(config *PublisherConfig) PublisherService {
	return &publisher{
		logger:   config.Logger,
		dao:      config.Dao,
		psClient: config.PSClient,
	}
}

func (p *publisher) ReadWriteKeys(ctx context.Context) error {
	var err error
	keyArr := utils.GetRandomStrings(100)
	processedIndex := make(map[int]bool)
	for {
		randomIndex := utils.GetRandomSequence(1, 100) - 1
		if len(processedIndex) == 100 {
			break
		}
		if processedIndex[randomIndex] {
			continue
		}
		processedIndex[randomIndex] = true
		if randomIndex < 30 {
			err = p.dao.InsertRecord(ctx, keyArr[randomIndex])
			if err != nil {
				break
			}
			err = p.psClient.PublishMessage(ctx, keyArr[randomIndex])
			if err != nil {
				break
			}
		} else {
			err = p.dao.GetRecord(ctx, keyArr[randomIndex])
			if err != nil {
				break
			}
		}

	}
	return err
}
