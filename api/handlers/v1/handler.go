package v1

import (
	"time"

	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"

	"musobaqa/farm-competition/internal/pkg/config"
	tokens "musobaqa/farm-competition/internal/pkg/token"

	appV "musobaqa/farm-competition/internal/usecase/app_version"
	// "musobaqa/farm-competition/internal/usecase/refresh_token"
)

type HandlerV1 struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	JwtHandler     tokens.JwtHandler
	AppVersion     appV.AppVersion
	Enforcer       *casbin.Enforcer
}

type HandlerV1Config struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	JwtHandler     tokens.JwtHandler
	AppVersion     appV.AppVersion
	Enforcer       *casbin.Enforcer
}

func New(c *HandlerV1Config) *HandlerV1 {
	return &HandlerV1{
		Config:         c.Config,
		Logger:         c.Logger,
		ContextTimeout: c.ContextTimeout,
		JwtHandler:     c.JwtHandler,
		AppVersion:     c.AppVersion,
		Enforcer:       c.Enforcer,
	}
}
