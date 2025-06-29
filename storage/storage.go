package storage

import (
	"context"
	"errors"
	"good_api_test/models"
)

var ErrNotFound = errors.New("not found")

type DB interface {
	GetGood(ctx context.Context, id, projectId int) (*models.Good, error)
	GetGoods(ctx context.Context, limit, offset int) ([]*models.Good, error)
	CreateGood(ctx context.Context, good *models.Good) (int, error)
	UpdateGood(ctx context.Context, good *models.Good) error
	DeleteGood(ctx context.Context, id, projectId int) error
	Reprioritize(ctx context.Context, id, projectId, newPriority int) ([]*models.Good, error)
	GetTotalGoodsCount(ctx context.Context) (int, error)
	GetRemovedGoodsCount(ctx context.Context) (int, error)
}
