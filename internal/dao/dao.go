package dao

import (
	"context"
	"log"
	"time"

	"github.com/aerospike/aerospike-client-go"
	"go.uber.org/zap"
)

type AeroDBService interface {
	InsertRecord(ctx context.Context, keyString string) error
	GetRecord(ctx context.Context, keyString string) error
}

type aeroDB struct {
	*aerospike.Client
	logger *zap.Logger
}

type AeroConfig struct {
	Host      string
	Port      int
	Logger    *zap.Logger
	QueueSize int
	LimitConn bool
	Timeout   int //in milliseconds
}

func GetDBConn(config AeroConfig) AeroDBService {
	clientPolicy := aerospike.NewClientPolicy()
	// clientPolicy.ConnectionQueueSize = config.QueueSize
	// clientPolicy.LimitConnectionsToQueueSize = config.LimitConn
	// clientPolicy.Timeout = time.Duration(config.Timeout) * time.Millisecond
	client, err := aerospike.NewClientWithPolicy(clientPolicy, config.Host, config.Port)
	if err != nil {
		log.Fatalf("client creation failed %s", err.Error())
	}
	return &aeroDB{
		logger: config.Logger,
		Client: client,
	}
}

func (ad *aeroDB) InsertRecord(ctx context.Context, keyString string) error {
	inputKey := keyString[:36]
	ad.logger.Info("inserting to aerodb", zap.String("inputKey", inputKey))
	key, err := aerospike.NewKey("test", "userData", inputKey)
	if err != nil {
		ad.logger.Error("failed to create key", zap.String("err", err.Error()))
		return err
	}
	dataBin := aerospike.NewBin("keyString",keyString)
	timeStampBin := aerospike.NewBin("createdAt",time.Now().String())

	writePolicy := aerospike.NewWritePolicy(0, 0)

	err = ad.PutBins(writePolicy, key, dataBin,timeStampBin)
	if err != nil {
		ad.logger.Error("failed to insert bin", zap.String("err", err.Error()))
		return err
	}
	ad.logger.Info("successfully inserted to aerodb", zap.String("keyString", keyString))
	return nil
}

func (ad *aeroDB) GetRecord(ctx context.Context, keyString string) error {
	ad.logger.Info("fetching from aerodb", zap.String("keyString", keyString))
	key, err := aerospike.NewKey("test", "userData", keyString)
	if err != nil {
		ad.logger.Error("failed to create key", zap.String("err", err.Error()))
		return err
	}

	rec, err := ad.Get(aerospike.NewPolicy(), key)
	if err != nil {
		if err.Error() == "Key not found" {
			ad.logger.Warn("key not found")
			return nil
		}
		ad.logger.Error("failed to get bin", zap.String("err", err.Error()))
		return err
	}
	ad.logger.Info("received bins", zap.Any("bins", rec.Bins))
	return nil
}
