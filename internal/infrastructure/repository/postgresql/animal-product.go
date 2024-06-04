package postgresql

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/spf13/cast"
	"musobaqa/farm-competition/internal/entity"
	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql/repo"
	"musobaqa/farm-competition/internal/pkg/postgres"
)

type animalProductRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewAnimalProduct(db *postgres.PostgresDB) repo.AnimalProduct {
	return &animalProductRepo{
		tableName: "animal_products",
		db:        db,
	}
}

func (ap *animalProductRepo) Create(ctx context.Context, animalProduct *entity.AnimalProductReq) (*entity.AnimalProductRes, error) {
	clauses := map[string]interface{}{
		"id":         animalProduct.ID,
		"animal_id":  animalProduct.AnimalID,
		"product_id": animalProduct.ProductID,
		"get_time":   animalProduct.GetTime,
		"capacity":   animalProduct.Capacity,
		"created_at": animalProduct.CreatedAt,
		"updated_at": animalProduct.UpdatedAt,
	}

	queryBuilder := ap.db.Sq.Builder.Insert(ap.tableName)
	queryBuilder = queryBuilder.SetMap(clauses)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	result, err := ap.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	if result.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	selectQueryBuilder := ap.db.Sq.Builder.Select(
		"a.id, " +
			"a.name, " +
			"a.category_name, " +
			"a.gender, " +
			"a.birth_day, " +
			"a.description, " +
			"a.genus, " +
			"a.weight, " +
			"a.is_health, " +
			"p.id, " +
			"p.name, " +
			"p.product_union, " +
			"p.description, " +
			"p.total_capacity, " +
			"ap.id, " +
			"ap.capacity, " +
			"ap.get_time")
	selectQueryBuilder = selectQueryBuilder.From(ap.tableName + " AS ap")
	selectQueryBuilder = selectQueryBuilder.Join("animals AS a ON a.id = ap.animal_id")
	selectQueryBuilder = selectQueryBuilder.Join("products AS p ON p.id = ap.product_id")
	selectQueryBuilder = selectQueryBuilder.Where("ap.deleted_at IS NULL")
	selectQueryBuilder = selectQueryBuilder.Where("a.deleted_at IS NULL")
	selectQueryBuilder = selectQueryBuilder.Where("p.deleted_at IS NULL")
	selectQueryBuilder = selectQueryBuilder.Where(ap.db.Sq.Equal("ap.id", animalProduct.ID))

	selectQuery, selectArgs, err := selectQueryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var (
		animalProductRes       entity.AnimalProductRes
		nullAnimalGenus        sql.NullString
		nullAnimalWeight       sql.NullInt64
		nullAnimalDescription  sql.NullString
		nullAnimalBirthday     sql.NullString
		nullProductDescription sql.NullString
		nullGetTime            sql.NullString
	)
	err = ap.db.QueryRow(ctx, selectQuery, selectArgs...).Scan(
		&animalProductRes.Animal.ID,
		&animalProductRes.Animal.Name,
		&animalProductRes.Animal.CategoryName,
		&animalProductRes.Animal.Gender,
		&nullAnimalBirthday,
		&nullAnimalDescription,
		&nullAnimalGenus,
		&nullAnimalWeight,
		&animalProductRes.Animal.IsHealth,
		&animalProductRes.Product.ID,
		&animalProductRes.Product.Name,
		&animalProductRes.Product.Union,
		&nullProductDescription,
		&animalProductRes.Product.TotalCapacity,
		&animalProductRes.ID,
		&animalProductRes.Capacity,
		&nullGetTime,
	)
	if err != nil {
		return nil, err
	}

	if nullAnimalBirthday.Valid {
		animalProductRes.Animal.BirthDay = nullAnimalBirthday.String
	}
	if nullAnimalGenus.Valid {
		animalProductRes.Animal.Genus = nullAnimalGenus.String
	}
	if nullAnimalWeight.Valid {
		animalProductRes.Animal.Weight = uint64(nullAnimalWeight.Int64)
	}
	if nullAnimalDescription.Valid {
		animalProductRes.Animal.Description = nullAnimalDescription.String
	}
	if nullProductDescription.Valid {
		animalProductRes.Product.Description = nullProductDescription.String
	}
	if nullGetTime.Valid {
		animalProductRes.GetTime = nullGetTime.String
	}

	return &animalProductRes, nil
}

func (ap *animalProductRepo) Update(ctx context.Context, animalProduct *entity.AnimalProductReq) (*entity.AnimalProductRes, error) {
	clauses := map[string]interface{}{
		"animal_id":  animalProduct.AnimalID,
		"product_id": animalProduct.ProductID,
		"get_time":   animalProduct.GetTime,
		"capacity":   animalProduct.Capacity,
		"updated_at": animalProduct.UpdatedAt,
	}

	queryBuilder := ap.db.Sq.Builder.Update(ap.tableName)
	queryBuilder = queryBuilder.SetMap(clauses)
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	queryBuilder = queryBuilder.Where(ap.db.Sq.Equal("id", animalProduct.ID))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	result, err := ap.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	if result.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	selectQueryBuilder := ap.db.Sq.Builder.Select(
		"a.id, " +
			"a.name, " +
			"a.category_name, " +
			"a.gender, " +
			"a.birth_day, " +
			"a.description, " +
			"a.genus, " +
			"a.weight, " +
			"a.is_health, " +
			"p.id, " +
			"p.name, " +
			"p.product_union, " +
			"p.description, " +
			"p.total_capacity, " +
			"ap.id, " +
			"ap.capacity, " +
			"ap.get_time")
	selectQueryBuilder = selectQueryBuilder.From(ap.tableName + " AS ap")
	selectQueryBuilder = selectQueryBuilder.Join("animals AS a ON a.id = ap.animal_id")
	selectQueryBuilder = selectQueryBuilder.Join("products AS p ON p.id = ap.product_id")
	selectQueryBuilder = selectQueryBuilder.Where("ap.deleted_at IS NULL")
	selectQueryBuilder = selectQueryBuilder.Where("a.deleted_at IS NULL")
	selectQueryBuilder = selectQueryBuilder.Where("p.deleted_at IS NULL")
	selectQueryBuilder = selectQueryBuilder.Where(ap.db.Sq.Equal("ap.id", animalProduct.ID))

	selectQuery, selectArgs, err := selectQueryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var (
		animalProductRes       entity.AnimalProductRes
		nullAnimalGenus        sql.NullString
		nullAnimalWeight       sql.NullInt64
		nullAnimalDescription  sql.NullString
		nullProductDescription sql.NullString
		nullAnimalBirthday     sql.NullString
		nullGetTime            sql.NullString
	)
	err = ap.db.QueryRow(ctx, selectQuery, selectArgs...).Scan(
		&animalProductRes.Animal.ID,
		&animalProductRes.Animal.Name,
		&animalProductRes.Animal.CategoryName,
		&animalProductRes.Animal.Gender,
		&nullAnimalBirthday,
		&nullAnimalDescription,
		&nullAnimalGenus,
		&nullAnimalWeight,
		&animalProductRes.Animal.IsHealth,
		&animalProductRes.Product.ID,
		&animalProductRes.Product.Name,
		&animalProductRes.Product.Union,
		&nullProductDescription,
		&animalProductRes.Product.TotalCapacity,
		&animalProductRes.ID,
		&animalProductRes.Capacity,
		&nullGetTime,
	)
	if err != nil {
		return nil, err
	}

	if nullAnimalBirthday.Valid {
		animalProductRes.Animal.BirthDay = nullAnimalBirthday.String
	}
	if nullAnimalGenus.Valid {
		animalProductRes.Animal.Genus = nullAnimalGenus.String
	}
	if nullAnimalWeight.Valid {
		animalProductRes.Animal.Weight = uint64(nullAnimalWeight.Int64)
	}
	if nullAnimalDescription.Valid {
		animalProductRes.Animal.Description = nullAnimalDescription.String
	}
	if nullProductDescription.Valid {
		animalProductRes.Product.Description = nullProductDescription.String
	}
	if nullGetTime.Valid {
		animalProductRes.GetTime = nullGetTime.String
	}

	return &animalProductRes, nil
}

func (ap *animalProductRepo) Delete(ctx context.Context, animalProductID string) error {
	queryBuilder := ap.db.Sq.Builder.Delete(ap.tableName)
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	queryBuilder = queryBuilder.Where(ap.db.Sq.Equal("id", animalProductID))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	result, err := ap.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (ap *animalProductRepo) Get(ctx context.Context, animalProductID string) (*entity.AnimalProductRes, error) {
	selectQueryBuilder := ap.db.Sq.Builder.Select(
		"a.id, " +
			"a.name, " +
			"a.category_name, " +
			"a.gender, " +
			"a.birth_day, " +
			"a.description, " +
			"a.genus, " +
			"a.weight, " +
			"a.is_health, " +
			"p.id, " +
			"p.name, " +
			"p.product_union, " +
			"p.description, " +
			"p.total_capacity, " +
			"ap.id, " +
			"ap.capacity, " +
			"ap.get_time")
	selectQueryBuilder = selectQueryBuilder.From(ap.tableName + " AS ap")
	selectQueryBuilder = selectQueryBuilder.Join("animals AS a ON a.id = ap.animal_id")
	selectQueryBuilder = selectQueryBuilder.Join("products AS p ON p.id = ap.product_id")
	selectQueryBuilder = selectQueryBuilder.Where("ap.deleted_at IS NULL")
	selectQueryBuilder = selectQueryBuilder.Where("a.deleted_at IS NULL")
	selectQueryBuilder = selectQueryBuilder.Where("p.deleted_at IS NULL")
	selectQueryBuilder = selectQueryBuilder.Where(ap.db.Sq.Equal("ap.id", animalProductID))

	selectQuery, selectArgs, err := selectQueryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var (
		animalProductRes       entity.AnimalProductRes
		nullAnimalGenus        sql.NullString
		nullAnimalWeight       sql.NullInt64
		nullAnimalDescription  sql.NullString
		nullProductDescription sql.NullString
		nullAnimalBirthday     sql.NullString
		nullGetTime            sql.NullString
	)
	err = ap.db.QueryRow(ctx, selectQuery, selectArgs...).Scan(
		&animalProductRes.Animal.ID,
		&animalProductRes.Animal.Name,
		&animalProductRes.Animal.CategoryName,
		&animalProductRes.Animal.Gender,
		&nullAnimalBirthday,
		&nullAnimalDescription,
		&nullAnimalGenus,
		&nullAnimalWeight,
		&animalProductRes.Animal.IsHealth,
		&animalProductRes.Product.ID,
		&animalProductRes.Product.Name,
		&animalProductRes.Product.Union,
		&nullProductDescription,
		&animalProductRes.Product.TotalCapacity,
		&animalProductRes.ID,
		&animalProductRes.Capacity,
		&nullGetTime,
	)
	if err != nil {
		return nil, err
	}

	if nullAnimalBirthday.Valid {
		animalProductRes.Animal.BirthDay = nullAnimalBirthday.String
	}
	if nullAnimalGenus.Valid {
		animalProductRes.Animal.Genus = nullAnimalGenus.String
	}
	if nullAnimalWeight.Valid {
		animalProductRes.Animal.Weight = uint64(nullAnimalWeight.Int64)
	}
	if nullAnimalDescription.Valid {
		animalProductRes.Animal.Description = nullAnimalDescription.String
	}
	if nullProductDescription.Valid {
		animalProductRes.Product.Description = nullProductDescription.String
	}
	if nullGetTime.Valid {
		animalProductRes.GetTime = nullGetTime.String
	}

	return &animalProductRes, nil
}

func (ap *animalProductRepo) List(ctx context.Context, page, limit uint64, params map[string]any) (*entity.ListAnimalProduct, error) {
	selectQueryBuilder := ap.db.Sq.Builder.Select(
		"a.id, " +
			"a.name, " +
			"a.category_name, " +
			"a.gender, " +
			"a.birth_day, " +
			"a.description, " +
			"a.genus, " +
			"a.weight, " +
			"a.is_health, " +
			"p.id, " +
			"p.name, " +
			"p.product_union, " +
			"p.description, " +
			"p.total_capacity, " +
			"ap.id, " +
			"ap.capacity, " +
			"ap.get_time")
	selectQueryBuilder = selectQueryBuilder.From(ap.tableName + " AS ap")
	selectQueryBuilder = selectQueryBuilder.Join("animals AS a ON a.id = ap.animal_id")
	selectQueryBuilder = selectQueryBuilder.Join("products AS p ON p.id = ap.product_id")
	selectQueryBuilder = selectQueryBuilder.Where("ap.deleted_at IS NULL")
	selectQueryBuilder = selectQueryBuilder.Where("a.deleted_at IS NULL")
	selectQueryBuilder = selectQueryBuilder.Where("p.deleted_at IS NULL")
	if params["get_time"] != "" {
		selectQueryBuilder = selectQueryBuilder.Where(ap.db.Sq.And(
			sq.GtOrEq{"ap.get_time": cast.ToString(params["get_time"]) + " 00:00:00______"},
			sq.LtOrEq{"ap.get_time": cast.ToString(params["get_time"]) + " 23:59:59______"},
		))
	}
	selectQueryBuilder = selectQueryBuilder.Limit(limit)
	selectQueryBuilder = selectQueryBuilder.Offset(limit * (page - 1))

	selectQuery, selectArgs, err := selectQueryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := ap.db.Query(ctx, selectQuery, selectArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var response entity.ListAnimalProduct
	for rows.Next() {
		var (
			animalProductRes       entity.AnimalProductRes
			nullAnimalGenus        sql.NullString
			nullAnimalWeight       sql.NullInt64
			nullAnimalDescription  sql.NullString
			nullProductDescription sql.NullString
			nullAnimalBirthday     sql.NullString
			nullGetTime            sql.NullString
		)
		err = rows.Scan(
			&animalProductRes.Animal.ID,
			&animalProductRes.Animal.Name,
			&animalProductRes.Animal.CategoryName,
			&animalProductRes.Animal.Gender,
			&nullAnimalBirthday,
			&nullAnimalDescription,
			&nullAnimalGenus,
			&nullAnimalWeight,
			&animalProductRes.Animal.IsHealth,
			&animalProductRes.Product.ID,
			&animalProductRes.Product.Name,
			&animalProductRes.Product.Union,
			&nullProductDescription,
			&animalProductRes.Product.TotalCapacity,
			&animalProductRes.ID,
			&animalProductRes.Capacity,
			&nullGetTime,
		)

		if err != nil {
			return nil, err
		}

		if nullAnimalBirthday.Valid {
			animalProductRes.Animal.BirthDay = nullAnimalBirthday.String
		}
		if nullAnimalGenus.Valid {
			animalProductRes.Animal.Genus = nullAnimalGenus.String
		}
		if nullAnimalWeight.Valid {
			animalProductRes.Animal.Weight = uint64(nullAnimalWeight.Int64)
		}
		if nullAnimalDescription.Valid {
			animalProductRes.Animal.Description = nullAnimalDescription.String
		}
		if nullProductDescription.Valid {
			animalProductRes.Product.Description = nullProductDescription.String
		}
		if nullGetTime.Valid {
			animalProductRes.GetTime = nullGetTime.String
		}

		response.AnimalProducts = append(response.AnimalProducts, &animalProductRes)
	}

	totalCountBuilder := ap.db.Sq.Builder.Select("COUNT(*)")
	totalCountBuilder = totalCountBuilder.From(ap.tableName + " AS ap")
	totalCountBuilder = totalCountBuilder.Join("animals AS a ON a.id = ap.animal_id")
	totalCountBuilder = totalCountBuilder.Join("products AS p ON p.id = ap.product_id")
	totalCountBuilder = totalCountBuilder.Where("ap.deleted_at IS NULL")
	totalCountBuilder = totalCountBuilder.Where("a.deleted_at IS NULL")
	totalCountBuilder = totalCountBuilder.Where("p.deleted_at IS NULL")
	if params["get_time"] != "" {
		totalCountBuilder = totalCountBuilder.Where(ap.db.Sq.And(
			sq.GtOrEq{"ap.get_time": cast.ToString(params["get_time"]) + " 00:00:00______"},
			sq.LtOrEq{"ap.get_time": cast.ToString(params["get_time"]) + " 23:59:59______"},
		))
	}

	totalQuery, totalArgs, err := totalCountBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var count = 0
	if err := ap.db.QueryRow(ctx, totalQuery, totalArgs...).Scan(&count); err != nil {
		return nil, err
	}
	response.TotalCount = uint64(count)

	return &response, nil
}

func (ap *animalProductRepo) ListAnimals(ctx context.Context, page, limit uint64, productID string) (*entity.AnimalsWithProduct, error) {
	queryBuilder := ap.db.Sq.Builder.Select("id, name, product_union, description, total_capacity")
	queryBuilder = queryBuilder.From("products")
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	queryBuilder = queryBuilder.Where(ap.db.Sq.Equal("id", productID))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var (
		response               entity.AnimalsWithProduct
		nullProductDescription sql.NullString
	)

	err = ap.db.QueryRow(ctx, query, args...).Scan(
		&response.Product.ID,
		&response.Product.Name,
		&response.Product.Union,
		&nullProductDescription,
		&response.Product.TotalCapacity,
	)
	if err != nil {
		return nil, err
	}

	if nullProductDescription.Valid {
		response.Product.Description = nullProductDescription.String
	}

	animalQueryBuilder := ap.db.Sq.Builder.Select("a.id, a.name, a.category_name, a.gender, a.birth_day, a.genus, a.weight, a.is_health, a.description, SUM(ap.capacity) AS total_category")
	animalQueryBuilder = animalQueryBuilder.From("animal_products AS ap")
	animalQueryBuilder = animalQueryBuilder.Join("animals AS a ON a.id = ap.animal_id")
	animalQueryBuilder = animalQueryBuilder.Where(ap.db.Sq.Equal("ap.product_id", productID))
	animalQueryBuilder = animalQueryBuilder.Where("ap.deleted_at IS NULL")
	animalQueryBuilder = animalQueryBuilder.Where("a.deleted_at IS NULL")
	animalQueryBuilder = animalQueryBuilder.GroupBy("a.id")
	animalQueryBuilder = animalQueryBuilder.OrderBy("a.total_category")
	animalQueryBuilder = animalQueryBuilder.Limit(limit)
	animalQueryBuilder = animalQueryBuilder.Offset(limit * (page - 1))

	animalQuery, animalArgs, err := animalQueryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := ap.db.Query(ctx, animalQuery, animalArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			animal                entity.Animal
			nullAnimalGenus       sql.NullString
			nullAnimalWeight      sql.NullInt64
			nullAnimalDescription sql.NullString
			NullAnimalBirthday    sql.NullString
			totosh                int64
		)
		err = rows.Scan(
			&animal.ID,
			&animal.Name,
			&animal.CategoryName,
			&animal.Gender,
			&NullAnimalBirthday,
			&nullAnimalGenus,
			&nullAnimalWeight,
			&nullAnimalDescription,
		)

		if err != nil {
			return nil, err
		}
		if nullAnimalGenus.Valid {
			animal.Genus = nullAnimalGenus.String
		}
		if nullAnimalWeight.Valid {
			animal.Weight = uint64(nullAnimalWeight.Int64)
		}
		if nullAnimalDescription.Valid {
			animal.Description = nullAnimalDescription.String
		}
		if nullProductDescription.Valid {
			animal.Description = nullProductDescription.String
		}

		response.Animals = append(response.Animals, &struct {
			ID            string
			Name          string
			CategoryName  string
			Gender        string
			BirthDay      string
			Genus         string
			Weight        uint64
			IsHealth      string
			Description   string
			TotalCapacity int64
		}{
			ID:            animal.ID,
			Name:          animal.Name,
			CategoryName:  animal.CategoryName,
			Gender:        animal.Gender,
			BirthDay:      animal.BirthDay,
			Genus:         animal.Genus,
			Weight:        animal.Weight,
			IsHealth:      animal.IsHealth,
			Description:   animal.Description,
			TotalCapacity: totosh,
		})
	}

	return &response, nil
}

func (ap *animalProductRepo) ListProducts(ctx context.Context, page, limit uint64, animalID string) (*entity.ProductsWithAnimal, error) {
	return nil, nil
}
