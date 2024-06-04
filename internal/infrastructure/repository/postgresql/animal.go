package postgresql

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v4"
	"musobaqa/farm-competition/internal/entity"
	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql/repo"
	"musobaqa/farm-competition/internal/pkg/postgres"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/spf13/cast"
)

type animalRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewAnimal(db *postgres.PostgresDB) repo.Animal {
	return &animalRepo{
		tableName: "animals",
		db:        db,
	}
}

func (a *animalRepo) Create(ctx context.Context, animal *entity.Animal) (*entity.Animal, error) {
	query := `
	INSERT INTO animals (
	    id,
	    name,
	    category_name,
	    gender,
	    birth_day,
	    genus,
	    weight,
	    description,
	    is_health,
	    created_at,
	    updated_at
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING
		id,
	    name,
		category_name,
		gender,
		birth_day,
		genus,
		weight,
		description,
		is_health
	`

	var (
		createdAnimal      entity.Animal
		sqlNullGenus       sql.NullString
		sqlNullDescription sql.NullString
		sqlNullBirthday    sql.NullString
	)

	err := a.db.QueryRow(ctx, query,
		animal.ID,
		animal.Name,
		animal.CategoryID,
		animal.Gender,
		animal.BirthDay,
		animal.Genus,
		animal.Weight,
		animal.Description,
		animal.IsHealth,
		animal.CreatedAt,
		animal.UpdatedAt,
	).Scan(
		&createdAnimal.ID,
		&createdAnimal.Name,
		&createdAnimal.CategoryID,
		&createdAnimal.Gender,
		&sqlNullBirthday,
		&sqlNullGenus,
		&createdAnimal.Weight,
		&sqlNullDescription,
		&createdAnimal.IsHealth,
	)

	if err != nil {
		return nil, err
	}
	if sqlNullGenus.Valid {
		createdAnimal.Genus = sqlNullGenus.String
	}
	if sqlNullDescription.Valid {
		createdAnimal.Description = sqlNullDescription.String
	}
	if sqlNullBirthday.Valid {
		createdAnimal.BirthDay = sqlNullBirthday.String
	}

	return &createdAnimal, nil
}

func (a *animalRepo) Update(ctx context.Context, animal *entity.Animal) (*entity.Animal, error) {
	query := `
	UPDATE
		animals
	SET
	    name = $1,
	    category_name= $2,
	    gender = $3,
	    birth_day = $4,
	    genus = $5,
	    weight = $6,
	    description = $7,
	    is_health = $8,
	    updated_at = $9
	WHERE
	    id = $10
		AND deleted_at IS NULL
	RETURNING
		id,
	    name,
		category_name,
		gender,
		birth_day,
		genus,
		weight,
		description,
		is_health
	`

	var (
		updatedAnimal      entity.Animal
		sqlNullGenus       sql.NullString
		sqlNullDescription sql.NullString
		sqlNullBirthday    sql.NullString
	)

	err := a.db.QueryRow(ctx, query,
		animal.Name,
		animal.CategoryID,
		animal.Gender,
		animal.BirthDay,
		animal.Genus,
		animal.Weight,
		animal.Description,
		animal.IsHealth,
		animal.UpdatedAt,
		animal.ID,
	).Scan(
		&updatedAnimal.ID,
		&updatedAnimal.Name,
		&updatedAnimal.CategoryID,
		&updatedAnimal.Gender,
		&sqlNullBirthday,
		&sqlNullGenus,
		&updatedAnimal.Weight,
		&sqlNullDescription,
		&updatedAnimal.IsHealth,
	)

	if err != nil {
		return nil, err
	}
	if sqlNullGenus.Valid {
		updatedAnimal.Genus = sqlNullGenus.String
	}
	if sqlNullDescription.Valid {
		updatedAnimal.Description = sqlNullDescription.String
	}
	if sqlNullBirthday.Valid {
		updatedAnimal.BirthDay = sqlNullBirthday.String
	}

	return &updatedAnimal, nil
}

func (a *animalRepo) Delete(ctx context.Context, animalID string) error {
	query := `UPDATE animals SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`

	result, err := a.db.Exec(ctx, query, time.Now().Format(time.RFC3339), animalID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (a *animalRepo) Get(ctx context.Context, animalID string) (*entity.Animal, error) {
	query := `
	SELECT
		id,
		name,
		category_name,
		gender,
		birth_day,
		genus,
		weight,
		description,
		is_health
	FROM
	    animals
	WHERE
	    id = $1
		AND deleted_at IS NULL
	`

	var (
		animal             entity.Animal
		sqlNullGenus       sql.NullString
		sqlNullDescription sql.NullString
		sqlNullBirthday    sql.NullString
	)

	err := a.db.QueryRow(ctx, query, animalID).Scan(
		&animal.ID,
		&animal.Name,
		&animal.CategoryID,
		&animal.Gender,
		&sqlNullBirthday,
		&sqlNullGenus,
		&animal.Weight,
		&sqlNullDescription,
		&animal.IsHealth,
	)
	if err != nil {
		return nil, err
	}

	if sqlNullGenus.Valid {
		animal.Genus = sqlNullGenus.String
	}
	if sqlNullDescription.Valid {
		animal.Description = sqlNullDescription.String
	}
	if sqlNullBirthday.Valid {
		animal.BirthDay = sqlNullBirthday.String
	}

	return &animal, nil
}

func (a *animalRepo) List(ctx context.Context, page, limit uint64, params map[string]any) (*entity.ListAnimal, error) {
	var (
		offset     = (page - 1) * limit
		animals    = entity.ListAnimal{}
		tenPercent = cast.ToInt(params["weight"]) / 10
		weightUp   = cast.ToInt(params["weight"]) + tenPercent
		weightDown = cast.ToInt(params["weight"]) - tenPercent
	)

	queryBuilder := a.db.Sq.Builder.Select("id, name, category_name, gender, birth_day, genus, weight, description, is_health")
	queryBuilder = queryBuilder.From(a.tableName)
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	queryBuilder = queryBuilder.Where(a.db.Sq.ILike("category_name", "%"+cast.ToString(params["category"])+"%"))
	queryBuilder = queryBuilder.Where(a.db.Sq.ILike("genus", "%"+cast.ToString(params["genus"])+"%"))
	queryBuilder = queryBuilder.Where(a.db.Sq.ILike("gender", "%"+cast.ToString(params["gender"])+"%"))
	if cast.ToInt(params["weight"]) != 0 {
		queryBuilder = queryBuilder.Where(a.db.Sq.And(
			sq.GtOrEq{"weight": weightDown},
			sq.LtOrEq{"weight": weightUp},
		))
		queryBuilder = queryBuilder.OrderBy("weight DESC")
	}

	if cast.ToString(params["is_health"]) != "" {
		queryBuilder = queryBuilder.Where(a.db.Sq.Equal("is_health", cast.ToString(params["is_health"])))
	}

	queryBuilder = queryBuilder.Limit(limit)
	queryBuilder = queryBuilder.Offset(offset)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := a.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			animal             entity.Animal
			sqlNullGenus       sql.NullString
			sqlNullDescription sql.NullString
			sqlNullBirthday    sql.NullString
		)
		err := rows.Scan(
			&animal.ID,
			&animal.Name,
			&animal.CategoryID,
			&animal.Gender,
			&sqlNullBirthday,
			&sqlNullGenus,
			&animal.Weight,
			&sqlNullDescription,
			&animal.IsHealth,
		)
		if err != nil {
			return nil, err
		}

		if sqlNullGenus.Valid {
			animal.Genus = sqlNullGenus.String
		}
		if sqlNullDescription.Valid {
			animal.Description = sqlNullDescription.String
		}
		if sqlNullBirthday.Valid {
			animal.BirthDay = sqlNullBirthday.String
		}

		animals.Animals = append(animals.Animals, &animal)
	}

	totalQueryBuilder := a.db.Sq.Builder.Select("COUNT(*)")
	totalQueryBuilder = totalQueryBuilder.From(a.tableName)
	totalQueryBuilder = totalQueryBuilder.Where("deleted_at IS NULL")
	totalQueryBuilder = totalQueryBuilder.Where(a.db.Sq.ILike("category_name", "%"+cast.ToString(params["category"])+"%"))
	totalQueryBuilder = totalQueryBuilder.Where(a.db.Sq.ILike("genus", "%"+cast.ToString(params["genus"])+"%"))
	totalQueryBuilder = totalQueryBuilder.Where(a.db.Sq.ILike("gender", "%"+cast.ToString(params["gender"])+"%"))
	if cast.ToString(params["is_health"]) != "" {
		totalQueryBuilder = totalQueryBuilder.Where(a.db.Sq.Equal("is_health", cast.ToString(params["is_health"])))
	}
	if cast.ToInt(params["weight"]) != 0 {
		totalQueryBuilder = totalQueryBuilder.Where(a.db.Sq.And(
			sq.GtOrEq{"weight": weightDown},
			sq.LtOrEq{"weight": weightUp},
		))
	}

	totalQuery, totalArgs, err := totalQueryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var count = 0
	if err := a.db.QueryRow(ctx, totalQuery, totalArgs...).Scan(&count); err != nil {
		return nil, err
	}
	animals.TotalCount = uint64(count)

	return &animals, nil
}
