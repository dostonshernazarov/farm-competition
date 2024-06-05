package postgresql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"musobaqa/farm-competition/internal/entity"
	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql/repo"
	"musobaqa/farm-competition/internal/pkg/postgres"
	"time"

	"github.com/jackc/pgx/v4"
)

type feedingRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewFeeding(db *postgres.PostgresDB) repo.Feeding {
	return &feedingRepo{
		tableName: "animal_given_eatables",
		db:        db,
	}
}

func (f *feedingRepo) Create(ctx context.Context, feeding *entity.Feeding) (*entity.FeedingRes, error) {
	clauses := map[string]interface{}{
		"id":          feeding.ID,
		"animal_id":   feeding.AnimalID,
		"eatables_id": feeding.EatablesID,
		"category":    feeding.Category,
		"day":         feeding.Day,
		"daily":       feeding.Daily,
		"created_at":  feeding.CreatedAt,
		"updated_at":  feeding.UpdatedAt,
	}

	var dayNull sql.NullString

	queryBuilder := f.db.Sq.Builder.Insert(f.tableName)
	queryBuilder = queryBuilder.SetMap(clauses)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	result, err := f.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	if result.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	selectFeedingBuilder := f.db.Sq.Builder.Select("id, animal_id, eatables_id, category, day, daily")
	selectFeedingBuilder = selectFeedingBuilder.From(f.tableName)
	selectFeedingBuilder = selectFeedingBuilder.Where("deleted_at IS NULL")
	selectFeedingBuilder = selectFeedingBuilder.Where(f.db.Sq.Equal("id", feeding.ID))

	selectQuery, selectArgs, err := selectFeedingBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var (
		dailyJson []byte
		response  entity.FeedingRes
	)
	err = f.db.QueryRow(ctx, selectQuery, selectArgs...).Scan(
		&response.ID,
		&response.AnimalID,
		&response.Eatables.ID,
		&response.Category,
		&dayNull,
		&dailyJson,
	)
	if err != nil {
		return nil, err
	}

	if dayNull.Valid {
		response.Day = dayNull.String
	}

	err = json.Unmarshal(dailyJson, &response.Daily)
	if err != nil {
		return nil, err
	}

	if feeding.Category == "food" {
		selectFoodBuilder := f.db.Sq.Builder.Select("id, name, capacity, description, product_union")
		selectFoodBuilder = selectFoodBuilder.From("foods")
		selectFoodBuilder = selectFoodBuilder.Where("deleted_at IS NULL")
		selectFoodBuilder = selectFoodBuilder.Where(f.db.Sq.Equal("id", feeding.EatablesID))

		selectFoodQuery, selectFoodArgs, err := selectFoodBuilder.ToSql()
		if err != nil {
			return nil, err
		}

		err = f.db.QueryRow(ctx, selectFoodQuery, selectFoodArgs...).Scan(
			&response.Eatables.ID,
			&response.Eatables.Name,
			&response.Eatables.Capacity,
			&response.Eatables.Description,
			&response.Eatables.Union,
		)
		if err != nil {
			return nil, err
		}
	} else if feeding.Category == "drug" {
		selectDrugBuilder := f.db.Sq.Builder.Select("id, name, status, capacity, description, product_union")
		selectDrugBuilder = selectDrugBuilder.From("drugs")
		selectDrugBuilder = selectDrugBuilder.Where("deleted_at IS NULL")
		selectDrugBuilder = selectDrugBuilder.Where(f.db.Sq.Equal("id", feeding.EatablesID))

		selectFoodQuery, selectFoodArgs, err := selectDrugBuilder.ToSql()
		if err != nil {
			return nil, err
		}

		err = f.db.QueryRow(ctx, selectFoodQuery, selectFoodArgs...).Scan(
			&response.Eatables.ID,
			&response.Eatables.Name,
			&response.Eatables.Status,
			&response.Eatables.Capacity,
			&response.Eatables.Description,
			&response.Eatables.Union,
		)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("unknown feeding category")
	}

	return &response, nil
}

func (f *feedingRepo) Update(ctx context.Context, feeding *entity.Feeding) (*entity.FeedingRes, error) {
	clauses := map[string]interface{}{
		"animal_id":   feeding.AnimalID,
		"eatables_id": feeding.EatablesID,
		"category":    feeding.Category,
		"day":         feeding.Day,
		"daily":       feeding.Daily,
		"updated_at":  feeding.UpdatedAt,
	}

	queryBuilder := f.db.Sq.Builder.Update(f.tableName)
	queryBuilder = queryBuilder.SetMap(clauses)
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	queryBuilder = queryBuilder.Where(f.db.Sq.Equal("id", feeding.ID))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	result, err := f.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	if result.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	selectFeedingBuilder := f.db.Sq.Builder.Select("id, animal_id, eatables_id, category, day, daily")
	selectFeedingBuilder = selectFeedingBuilder.From(f.tableName)
	selectFeedingBuilder = selectFeedingBuilder.Where("deleted_at IS NULL")
	selectFeedingBuilder = selectFeedingBuilder.Where(f.db.Sq.Equal("id", feeding.ID))

	selectQuery, selectArgs, err := selectFeedingBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var (
		dailyJson []byte
		response  entity.FeedingRes
		dayNull sql.NullString
	)
	err = f.db.QueryRow(ctx, selectQuery, selectArgs...).Scan(
		&response.ID,
		&response.AnimalID,
		&response.Eatables.ID,
		&response.Category,
		&dayNull,
		&dailyJson,
	)
	if err != nil {
		return nil, err
	}

	if dayNull.Valid {
		response.Day = dayNull.String
	}

	err = json.Unmarshal(dailyJson, &response.Daily)
	if err != nil {
		return nil, err
	}

	if feeding.Category == "food" {
		selectFoodBuilder := f.db.Sq.Builder.Select("id, name, capacity, description, product_union")
		selectFoodBuilder = selectFoodBuilder.From("foods")
		selectFoodBuilder = selectFoodBuilder.Where("deleted_at IS NULL")
		selectFoodBuilder = selectFoodBuilder.Where(f.db.Sq.Equal("id", feeding.EatablesID))

		selectFoodQuery, selectFoodArgs, err := selectFoodBuilder.ToSql()
		if err != nil {
			return nil, err
		}

		err = f.db.QueryRow(ctx, selectFoodQuery, selectFoodArgs...).Scan(
			&response.Eatables.ID,
			&response.Eatables.Name,
			&response.Eatables.Capacity,
			&response.Eatables.Description,
			&response.Eatables.Union,
		)
		if err != nil {
			return nil, err
		}
	} else if feeding.Category == "drug" {
		selectDrugBuilder := f.db.Sq.Builder.Select("id, name, status, capacity, description, product_union")
		selectDrugBuilder = selectDrugBuilder.From("drugs")
		selectDrugBuilder = selectDrugBuilder.Where("deleted_at IS NULL")
		selectDrugBuilder = selectDrugBuilder.Where(f.db.Sq.Equal("id", feeding.EatablesID))

		selectFoodQuery, selectFoodArgs, err := selectDrugBuilder.ToSql()
		if err != nil {
			return nil, err
		}

		err = f.db.QueryRow(ctx, selectFoodQuery, selectFoodArgs...).Scan(
			&response.Eatables.ID,
			&response.Eatables.Name,
			&response.Eatables.Status,
			&response.Eatables.Capacity,
			&response.Eatables.Description,
			&response.Eatables.Union,
		)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("unknown feeding category")
	}

	return &response, nil
}

func (f *feedingRepo) Delete(ctx context.Context, feedingID string) error {
	clauses := map[string]interface{}{
		"deleted_at": time.Now().Format(time.RFC3339),
	}

	queryBuilder := f.db.Sq.Builder.Update(f.tableName)
	queryBuilder = queryBuilder.SetMap(clauses)
	queryBuilder = queryBuilder.Where("deleted_at is null")
	queryBuilder = queryBuilder.Where(f.db.Sq.Equal("id", feedingID))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	result, err := f.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
