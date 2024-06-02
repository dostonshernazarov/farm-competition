package postgresql

import (
	"musobaqa/farm-competition/internal/pkg/postgres"
	"musobaqa/farm-competition/internal/usecase/refresh_token"
)

func NewRefreshTokenRepo(db *postgres.PostgresDB) refresh_token.RefreshTokenRepo {
	return nil
}
