package storage

import (
	"context"
	"errors"
	"github.com/eqkez0r/gophermart/internal/storage/postgres"
	"go.uber.org/zap"
)

var (
	ErrUnknownStorageType = errors.New("Unsupported storage type")
)

func NewStorage(
	ctx context.Context,
	logger *zap.SugaredLogger,
	storagetype string,
	settings ...string) (Storage, error) {
	switch storagetype {
	case "postgresql":
		{
			return postgres.New(ctx, logger, settings[0])
		}
	default:
		return nil, ErrUnknownStorageType
	}
}
