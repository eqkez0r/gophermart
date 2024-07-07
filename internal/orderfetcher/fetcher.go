package orderfetcher

import (
	"context"
	"github.com/eqkez0r/gophermart/internal/storage"
	"go.uber.org/zap"
)

type OrdersProvider interface {
	GetUnfinishedOrders(ctx context.Context) error
}

type OrderFetcher struct {
	storage    storage.Storage
	logger     *zap.SugaredLogger
	accrualuri string
}

func New(
	s storage.Storage,
	logger *zap.SugaredLogger,
	accrualuri string,
) *OrderFetcher {

	return &OrderFetcher{
		storage:    s,
		logger:     logger,
		accrualuri: accrualuri,
	}
}

func (or *OrderFetcher) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:

		}
	}
}

func (or *OrderFetcher) getOrders(ctx context.Context) error {
	return nil
}
