package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"good_api_test/models"
)

type PostgresDB struct {
	pool *pgxpool.Pool
}

func NewPostgresDB(connString string) (*PostgresDB, error) {
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	return &PostgresDB{pool},
		nil
}

func (db *PostgresDB) GetGood(ctx context.Context, id, projectId int) (*models.Good, error) {
	good := &models.Good{}
	err := db.pool.QueryRow(ctx,
		"SELECT id, project_id, name, description, priority, removed, created_at FROM goods WHERE id=$1 AND project_id=$2",
		id, projectId).Scan(&good.Id, &good.ProjectId, &good.Name, &good.Description, &good.Priority, &good.Removed, &good.CreatedAt)
	if err != nil {
		return nil, err
	}
	return good, nil
}

func (db *PostgresDB) GetGoods(ctx context.Context, limit, offset int) ([]*models.Good, error) {
	rows, err := db.pool.Query(ctx,
		"SELECT id, project_id, name, description, priority, removed, created_at FROM goods ORDER BY priority LIMIT $1 OFFSET $2",
		limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goods []*models.Good
	for rows.Next() {
		good := &models.Good{}
		err := rows.Scan(&good.Id, &good.ProjectId, &good.Name, &good.Description, &good.Priority, &good.Removed, &good.CreatedAt)
		if err != nil {
			return nil, err
		}
		goods = append(goods, good)
	}

	return goods, nil
}

func (db *PostgresDB) CreateGood(ctx context.Context, good *models.Good) (int, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "SELECT 1 FROM goods FOR UPDATE")
	if err != nil {
		return 0, err
	}

	var maxPriority int
	err = tx.QueryRow(ctx, "SELECT COALESCE(MAX(priority), 0) FROM goods").Scan(&maxPriority)
	if err != nil {
		return 0, err
	}

	good.Priority = maxPriority + 1

	var id int
	err = tx.QueryRow(ctx,
		"INSERT INTO goods (project_id, name, description, priority) VALUES ($1, $2, $3, $4) RETURNING id",
		good.ProjectId, good.Name, good.Description, good.Priority).Scan(&id)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return id, nil
}

func (db *PostgresDB) UpdateGood(ctx context.Context, good *models.Good) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	res, err := tx.Exec(ctx, "SELECT 1 FROM goods WHERE id=$1 AND project_id=$2 FOR UPDATE", good.Id, good.ProjectId)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	res, err = tx.Exec(ctx,
		"UPDATE goods SET name=$1, description=$2 WHERE id=$3 AND project_id=$4",
		good.Name, good.Description, good.Id, good.ProjectId)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	return tx.Commit(ctx)
}

func (db *PostgresDB) DeleteGood(ctx context.Context, id, projectId int) error {
	res, err := db.pool.Exec(ctx, "UPDATE goods SET removed=true WHERE id=$1 AND project_id=$2", id, projectId)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (db *PostgresDB) Reprioritize(ctx context.Context, id, projectId, newPriority int) ([]*models.Good, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var currentPriority int
	err = tx.QueryRow(ctx, "SELECT priority FROM goods WHERE id=$1 AND project_id=$2", id, projectId).Scan(&currentPriority)
	if err != nil {
		return nil, err
	}

	if newPriority == currentPriority {
		return nil, nil
	}

	rows, err := tx.Query(ctx,
		`UPDATE goods
		 SET priority = CASE
			 WHEN id = $1 AND project_id = $2 THEN $3
			 WHEN $4::INTEGER > $5::INTEGER AND priority > $5::INTEGER AND priority <= $4::INTEGER THEN priority - 1
			 WHEN $4::INTEGER < $5::INTEGER AND priority >= $4::INTEGER AND priority < $5::INTEGER THEN priority + 1
			 ELSE priority
		 END
		 WHERE (id = $1 AND project_id = $2) OR
			   ($4::INTEGER > $5::INTEGER AND priority > $5::INTEGER AND priority <= $4::INTEGER) OR
			   ($4::INTEGER < $5::INTEGER AND priority >= $4::INTEGER AND priority < $5::INTEGER)
		 RETURNING id, project_id, name, description, priority, removed, created_at`,
		id, projectId, newPriority, newPriority, currentPriority)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goods []*models.Good
	for rows.Next() {
		good := &models.Good{}
		err := rows.Scan(&good.Id, &good.ProjectId, &good.Name, &good.Description, &good.Priority, &good.Removed, &good.CreatedAt)
		if err != nil {
			return nil, err
		}
		goods = append(goods, good)
	}

	return goods, tx.Commit(ctx)
}

func (db *PostgresDB) GetTotalGoodsCount(ctx context.Context) (int, error) {
	var count int
	err := db.pool.QueryRow(ctx, "SELECT COUNT(*) FROM goods").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (db *PostgresDB) GetRemovedGoodsCount(ctx context.Context) (int, error) {
	var count int
	err := db.pool.QueryRow(ctx, "SELECT COUNT(*) FROM goods WHERE removed = true").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
