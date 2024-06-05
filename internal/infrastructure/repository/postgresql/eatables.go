package postgresql

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v4"
	"musobaqa/farm-competition/internal/entity"
	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql/repo"
	"musobaqa/farm-competition/internal/pkg/postgres"
	"time"
)

type eatableRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewEatable(db *postgres.PostgresDB) repo.Eatable {
	return &eatableRepo{
		tableName: "animal_eatable_info",
		db:        db,
	}
}

func (e *eatableRepo) Create(ctx context.Context, eatable *entity.Eatables) (*entity.EatablesRes, error) {
	dailyJson, err := json.Marshal(eatable.Daily)
	if err != nil {
		return nil, err
	}

	clauses := map[string]interface{}{
		"id":          eatable.ID,
		"animal_id":   eatable.AnimalID,
		"eatables_id": eatable.EatableID,
		"category":    eatable.Category,
		"daily":       dailyJson,
		"created_at":  eatable.CreatedAt,
		"updated_at":  eatable.UpdatedAt,
	}

	queryBuilder := e.db.Sq.Builder.Insert(e.tableName)
	queryBuilder = queryBuilder.SetMap(clauses)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	result, err := e.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	if result.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	selectEatableBuilder := e.db.Sq.Builder.Select("id, animal_id, eatables_id, category, daily")
	selectEatableBuilder = selectEatableBuilder.From(e.tableName)
	selectEatableBuilder = selectEatableBuilder.Where("deleted_at IS NULL")
	selectEatableBuilder = selectEatableBuilder.Where(e.db.Sq.Equal("id", eatable.ID))

	selectEatableQuery, selectEatableArgs, err := selectEatableBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var (
		dailyByte []byte
		response  entity.EatablesRes
	)
	err = e.db.QueryRow(ctx, selectEatableQuery, selectEatableArgs...).Scan(
		&response.ID,
		&response.AnimalID,
		&response.Eatable.ID,
		&response.Category,
		&dailyByte,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(dailyByte, &response.Daily)
	if err != nil {
		return nil, err
	}

	if eatable.Category == "food" {
		selectFoodBuilder := e.db.Sq.Builder.Select("id, name, capacity, description, product_union")
		selectFoodBuilder = selectFoodBuilder.From("foods")
		selectFoodBuilder = selectFoodBuilder.Where("deleted_at IS NULL")
		selectFoodBuilder = selectFoodBuilder.Where(e.db.Sq.Equal("id", eatable.EatableID))

		selectFoodQuery, selectFoodArgs, err := selectFoodBuilder.ToSql()
		if err != nil {
			return nil, err
		}

		err = e.db.QueryRow(ctx, selectFoodQuery, selectFoodArgs...).Scan(
			&response.Eatable.ID,
			&response.Eatable.Name,
			&response.Eatable.Capacity,
			&response.Eatable.Description,
			&response.Eatable.Union,
		)
		if err != nil {
			return nil, err
		}
	} else if eatable.Category == "drug" {
		selectDrugBuilder := e.db.Sq.Builder.Select("id, name, status, capacity, description, product_union")
		selectDrugBuilder = selectDrugBuilder.From("drugs")
		selectDrugBuilder = selectDrugBuilder.Where("deleted_at IS NULL")
		selectDrugBuilder = selectDrugBuilder.Where(e.db.Sq.Equal("id", eatable.EatableID))

		selectFoodQuery, selectFoodArgs, err := selectDrugBuilder.ToSql()
		if err != nil {
			return nil, err
		}

		err = e.db.QueryRow(ctx, selectFoodQuery, selectFoodArgs...).Scan(
			&response.Eatable.ID,
			&response.Eatable.Name,
			&response.Eatable.Status,
			&response.Eatable.Capacity,
			&response.Eatable.Description,
			&response.Eatable.Union,
		)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("unknown eatable category")
	}

	return &response, nil
}

func (e *eatableRepo) Update(ctx context.Context, eatable *entity.Eatables) (*entity.EatablesRes, error) {
	dailyJson, err := json.Marshal(eatable.Daily)
	if err != nil {
		return nil, err
	}
	clauses := map[string]interface{}{
		"animal_id":   eatable.AnimalID,
		"eatables_id": eatable.EatableID,
		"category":    eatable.Category,
		"daily":       dailyJson,
		"updated_at":  eatable.UpdatedAt,
	}

	queryBuilder := e.db.Sq.Builder.Update(e.tableName)
	queryBuilder = queryBuilder.SetMap(clauses)
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	queryBuilder = queryBuilder.Where(e.db.Sq.Equal("id", eatable.ID))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	result, err := e.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	if result.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	selectEatableBuilder := e.db.Sq.Builder.Select("id, animal_id, eatables_id, category, daily")
	selectEatableBuilder = selectEatableBuilder.From(e.tableName)
	selectEatableBuilder = selectEatableBuilder.Where("deleted_at IS NULL")
	selectEatableBuilder = selectEatableBuilder.Where(e.db.Sq.Equal("id", eatable.ID))

	selectEatableQuery, selectEatableArgs, err := selectEatableBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var (
		dailyByte []byte
		response  entity.EatablesRes
	)
	err = e.db.QueryRow(ctx, selectEatableQuery, selectEatableArgs...).Scan(
		&response.ID,
		&response.AnimalID,
		&response.Eatable.ID,
		&response.Category,
		&dailyByte,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(dailyByte, &response.Daily)
	if err != nil {
		return nil, err
	}

	if eatable.Category == "food" {
		selectFoodBuilder := e.db.Sq.Builder.Select("id, name, capacity, description, product_union")
		selectFoodBuilder = selectFoodBuilder.From("foods")
		selectFoodBuilder = selectFoodBuilder.Where("deleted_at IS NULL")
		selectFoodBuilder = selectFoodBuilder.Where(e.db.Sq.Equal("id", eatable.EatableID))

		selectFoodQuery, selectFoodArgs, err := selectFoodBuilder.ToSql()
		if err != nil {
			return nil, err
		}

		err = e.db.QueryRow(ctx, selectFoodQuery, selectFoodArgs...).Scan(
			&response.Eatable.ID,
			&response.Eatable.Name,
			&response.Eatable.Capacity,
			&response.Eatable.Description,
			&response.Eatable.Union,
		)
		if err != nil {
			return nil, err
		}
	} else if eatable.Category == "drug" {
		selectDrugBuilder := e.db.Sq.Builder.Select("id, name, status, capacity, description, product_union")
		selectDrugBuilder = selectDrugBuilder.From("drugs")
		selectDrugBuilder = selectDrugBuilder.Where("deleted_at IS NULL")
		selectDrugBuilder = selectDrugBuilder.Where(e.db.Sq.Equal("id", eatable.EatableID))

		selectFoodQuery, selectFoodArgs, err := selectDrugBuilder.ToSql()
		if err != nil {
			return nil, err
		}

		err = e.db.QueryRow(ctx, selectFoodQuery, selectFoodArgs...).Scan(
			&response.Eatable.ID,
			&response.Eatable.Name,
			&response.Eatable.Status,
			&response.Eatable.Capacity,
			&response.Eatable.Description,
			&response.Eatable.Union,
		)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("unknown eatable category")
	}

	return &response, nil
}

func (e *eatableRepo) Delete(ctx context.Context, eatablesID string) error {
	queryBuilder := e.db.Sq.Builder.Update(e.tableName)
	queryBuilder = queryBuilder.SetMap(map[string]interface{}{
		"deleted_at": time.Now().Format(time.RFC3339),
	})
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	queryBuilder = queryBuilder.Where(e.db.Sq.Equal("id", eatablesID))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}
	result, err := e.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (e *eatableRepo) GetFoods(ctx context.Context, page, limit uint64, animalID string) (*entity.ListFoodEatables, error) {
	queryBuilder := e.db.Sq.Builder.Select(
		"e.id, " +
			"e.animal_id, " +
			"e.daily, " +
			"f.id, " +
			"f.name, " +
			"f.capacity, " +
			"f.description, " +
			"f.product_union")
	queryBuilder = queryBuilder.From(e.tableName + " AS e")
	queryBuilder = queryBuilder.Join("foods AS f ON f.id = e.eatables_id")
	queryBuilder = queryBuilder.Where("f.deleted_at IS NULL")
	queryBuilder = queryBuilder.Where("e.deleted_at IS NULL")
	queryBuilder = queryBuilder.Where(e.db.Sq.Equal("e.animal_id", animalID))
	queryBuilder = queryBuilder.Limit(limit)
	queryBuilder = queryBuilder.Offset(limit * (page - 1))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := e.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var response entity.ListFoodEatables
	for rows.Next() {
		var (
			eatableDailyByte []byte
			eatables         entity.EatablesFoodRes
		)
		err = rows.Scan(
			&eatables.ID,
			&eatables.AnimalID,
			&eatableDailyByte,
			&eatables.Food.ID,
			&eatables.Food.Name,
			&eatables.Food.Capacity,
			&eatables.Food.Description,
			&eatables.Food.Union,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(eatableDailyByte, &eatables.Daily)
		if err != nil {
			return nil, err
		}

		response.Eatables = append(response.Eatables, &eatables)
	}

	countBuilder := e.db.Sq.Builder.Select("COUNT(*)")
	countBuilder = countBuilder.From(e.tableName + " AS e")
	countBuilder = countBuilder.Join("foods AS f ON f.id = e.eatables_id")
	countBuilder = countBuilder.Where("f.deleted_at IS NULL")
	countBuilder = countBuilder.Where("e.deleted_at IS NULL")
	countBuilder = countBuilder.Where(e.db.Sq.Equal("e.animal_id", animalID))

	countQuery, countArgs, err := countBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var count int64
	err = e.db.QueryRow(ctx, countQuery, countArgs...).Scan(&count)
	if err != nil {
		return nil, err
	}
	response.TotalCount = uint64(count)

	return &response, nil
}

func (e *eatableRepo) GetDrugs(ctx context.Context, page, limit uint64, animalID string) (*entity.ListDrugEatables, error) {
	queryBuilder := e.db.Sq.Builder.Select(
		"e.id, " +
			"e.animal_id, " +
			"e.daily, " +
			"d.id, " +
			"d.name, " +
			"d.status, " +
			"d.capacity, " +
			"d.description, " +
			"d.product_union")
	queryBuilder = queryBuilder.From(e.tableName + " AS e")
	queryBuilder = queryBuilder.Join("drugs AS d ON d.id = e.eatables_id")
	queryBuilder = queryBuilder.Where("d.deleted_at IS NULL")
	queryBuilder = queryBuilder.Where("e.deleted_at IS NULL")
	queryBuilder = queryBuilder.Where(e.db.Sq.Equal("e.animal_id", animalID))
	queryBuilder = queryBuilder.Limit(limit)
	queryBuilder = queryBuilder.Offset(limit * (page - 1))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := e.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var response entity.ListDrugEatables
	for rows.Next() {
		var (
			eatableDailyByte []byte
			eatables         entity.EatablesDrugRes
		)
		err = rows.Scan(
			&eatables.ID,
			&eatables.AnimalID,
			&eatableDailyByte,
			&eatables.Drug.ID,
			&eatables.Drug.Name,
			&eatables.Drug.Status,
			&eatables.Drug.Capacity,
			&eatables.Drug.Description,
			&eatables.Drug.Union,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(eatableDailyByte, &eatables.Daily)
		if err != nil {
			return nil, err
		}

		response.Eatables = append(response.Eatables, &eatables)
	}

	countBuilder := e.db.Sq.Builder.Select("COUNT(*)")
	countBuilder = countBuilder.From(e.tableName + " AS e")
	countBuilder = countBuilder.Join("drugs AS d ON d.id = e.eatables_id")
	countBuilder = countBuilder.Where("d.deleted_at IS NULL")
	countBuilder = countBuilder.Where("e.deleted_at IS NULL")
	countBuilder = countBuilder.Where(e.db.Sq.Equal("e.animal_id", animalID))

	countQuery, countArgs, err := countBuilder.ToSql()
	if err != nil {
		return nil, err
	}
	var count int64
	err = e.db.QueryRow(ctx, countQuery, countArgs...).Scan(&count)
	if err != nil {
		return nil, err
	}
	response.TotalCount = uint64(count)

	return &response, nil
}
