package services

import (
	"context"
	"errors"
	clientMocks "poc-publisher/internal/client/mocks"
	daoMocks "poc-publisher/internal/dao/mocks"
	"poc-publisher/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_ReadWriteKeys(t *testing.T) {
	tests := []struct {
		name        string
		index       int
		insertErr   error
		publishErr  error
		getErr      error
		errExpected bool
	}{
		{
			name:  "Write Success ",
			index: 1,
		},
		{
			name: "Read Success",
			index: 300,
		},
		{
			name: "Insert Failure",
			index: 1,
			insertErr: errors.New("insert error"),
			errExpected: true,
		},
		{
			name: "Publish Failure",
			index: 1,
			publishErr: errors.New("publish error"),
			errExpected: true,
		},
		{
			name: "Get Failure",
			index: 300,
			getErr: errors.New("get error"),
			errExpected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			keyArr := utils.StaticStrings
			logger, _ := zap.NewProduction()
			dao := &daoMocks.AeroDBService{}
			client := &clientMocks.PubSubClient{}
			dao.On("InsertRecord", ctx, keyArr[tt.index]).Return(tt.insertErr)
			client.On("PublishMessage", ctx, keyArr[tt.index]).Return(tt.publishErr)
			dao.On("GetRecord", ctx, keyArr[tt.index]).Return(tt.getErr)
			publisherSrv := GetPublisherService(&PublisherConfig{
				Logger:   logger,
				Dao:      dao,
				PSClient: client,
			})
			err := publisherSrv.ReadWriteKeys(ctx, tt.index)
			assert.Equal(t, tt.errExpected, err != nil)
		})
	}

}
