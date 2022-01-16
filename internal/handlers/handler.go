package handlers

import (
	"context"
	"net/http"
	"poc-publisher/internal/literals"
	"poc-publisher/internal/services"
	"poc-publisher/internal/utils"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type RouterConfig struct {
	Service services.PublisherService
	Logger  *zap.Logger
}

func InitRouter(ctx context.Context, config RouterConfig) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/aerospike-testing", getSetFromAerospike(ctx, config)).Methods(http.MethodGet)
	return router
}

func getSetFromAerospike(ctx context.Context, config RouterConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		config.Logger.Info("routed to getSetFromAerospike....")
		w.Header().Set("Content-Type", "text/plain")
		randSeq := utils.GetRandomSequence(1, 1000) - 1
		err := config.Service.ReadWriteKeys(ctx, randSeq)
		if err != nil {
			if err.Error() == literals.KEY_NOT_FOUND {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(err.Error()))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		if randSeq < 300 {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("db write: OK"))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("db read: OK"))
		}

	}
}
