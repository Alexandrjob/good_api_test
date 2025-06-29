package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"good_api_test/broker"
	"good_api_test/cache"
	"good_api_test/models"
	"good_api_test/storage"
)

var ErrGoodNotFound = errors.New("good not found")

type Service struct {
	storage      storage.DB
	cache        cache.Cache
	broker       broker.Broker
	publishQueue chan *models.Good
}

func New(storage storage.DB, cache cache.Cache, broker broker.Broker) *Service {
	s := &Service{storage, cache, broker, make(chan *models.Good, 100)}
	go s.startPublisher()
	return s
}

type GoodsListResponse struct {
	Meta  Meta           `json:"meta"`
	Goods []*models.Good `json:"goods"`
}

type Meta struct {
	Total   int `json:"total"`
	Removed int `json:"removed"`
	Limit   int `json:"limit"`
	Offset  int `json:"offset"`
}

func (s *Service) Get(ctx context.Context, id, projectId int) (*models.Good, error) {
	cacheKey := fmt.Sprintf("good:%d:%d", id, projectId)
	good, err := s.cache.Get(ctx, cacheKey)
	if err == nil {
		return good, nil
	}

	good, err = s.storage.GetGood(ctx, id, projectId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrGoodNotFound
		}
		return nil, err
	}

	s.cache.Set(ctx, cacheKey, good, time.Minute)
	return good, nil
}

func (s *Service) GetAll(ctx context.Context, limit, offset int) (*GoodsListResponse, error) {
	goods, err := s.storage.GetGoods(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := s.storage.GetTotalGoodsCount(ctx)
	if err != nil {
		return nil, err
	}

	removed, err := s.storage.GetRemovedGoodsCount(ctx)
	if err != nil {
		return nil, err
	}

	return &GoodsListResponse{
		Goods: goods,
		Meta: Meta{
			Total:   total,
			Removed: removed,
			Limit:   limit,
			Offset:  offset,
		},
	}, nil
}

func (s *Service) Create(ctx context.Context, good *models.Good) (int, error) {
	id, err := s.storage.CreateGood(ctx, good)
	if err != nil {
		return 0, err
	}
	good.Id = id
	s.publishQueue <- good
	return id, nil
}

func (s *Service) Update(ctx context.Context, id int, good *models.Good) error {
	if err := s.storage.UpdateGood(ctx, good); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrGoodNotFound
		}
		return err
	}

	s.cache.Delete(ctx, fmt.Sprintf("good:%d:%d", id, good.ProjectId))
	fetchedGood, err := s.storage.GetGood(ctx, id, good.ProjectId)
	if err == nil {
		s.publishQueue <- fetchedGood
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, id, projectId int) error {
	if err := s.storage.DeleteGood(ctx, id, projectId); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrGoodNotFound
		}
		return err
	}
	s.cache.Delete(ctx, fmt.Sprintf("good:%d:%d", id, projectId))
	deletedGood := &models.Good{Id: id, ProjectId: projectId, Removed: true}
	s.publishQueue <- deletedGood
	return nil
}

func (s *Service) Reprioritize(ctx context.Context, id, projectId, newPriority int) ([]*models.Good, error) {
	goods, err := s.storage.Reprioritize(ctx, id, projectId, newPriority)
	if err != nil {
		return nil, err
	}

	for _, good := range goods {
		s.cache.Delete(ctx, fmt.Sprintf("good:%d:%d", good.Id, good.ProjectId))
		s.publishQueue <- good
	}
	return goods, nil
}

func (s *Service) startPublisher() {
	for good := range s.publishQueue {
		s.broker.Publish(context.Background(), good)
	}
}
