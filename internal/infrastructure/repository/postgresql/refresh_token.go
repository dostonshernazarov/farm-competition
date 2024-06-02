package postgresql

import (
	"musobaqa/api-service/internal/pkg/postgres"
	"musobaqa/api-service/internal/usecase/refresh_token"
)

func NewRefreshTokenRepo(db *postgres.PostgresDB) refresh_token.RefreshTokenRepo {
	return nil
}
