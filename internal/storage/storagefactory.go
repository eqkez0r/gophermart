package storage

import (
	"context"
	"github.com/eqkez0r/gophermart/internal/storage/postgres"
	"go.uber.org/zap"
)

func NewStorage(
	ctx context.Context,
	logger *zap.SugaredLogger,
	storagetype string,
	settings ...string) Storage {
	switch storagetype {
	case "postgrsql":
		{
			return postgres.New()
		}
	default:
		return nil
	}
}
