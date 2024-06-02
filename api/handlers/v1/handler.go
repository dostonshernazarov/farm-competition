package v1

import (
	"time"

	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"

	grpcClients "musobaqa/api-service/internal/infrastructure/grpc_service_client"
	"musobaqa/api-service/internal/pkg/config"
	tokens "musobaqa/api-service/internal/pkg/token"

	appV "musobaqa/api-service/internal/usecase/app_version"
	"musobaqa/api-service/internal/usecase/event"
	// "musobaqa/api-service/internal/usecase/refresh_token"
)

type HandlerV1 struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	JwtHandler     tokens.JwtHandler
	Service        grpcClients.ServiceClient
	AppVersion     appV.AppVersion
	BrokerProducer event.BrokerProducer
	Enforcer       *casbin.Enforcer
}

type HandlerV1Config struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	JwtHandler     tokens.JwtHandler
	Service        grpcClients.ServiceClient
	AppVersion     appV.AppVersion
	BrokerProducer event.BrokerProducer
	Enforcer       *casbin.Enforcer
}

func New(c *HandlerV1Config) *HandlerV1 {
	return &HandlerV1{
		Config:         c.Config,
		Logger:         c.Logger,
		ContextTimeout: c.ContextTimeout,
		Service:        c.Service,
		JwtHandler: c.JwtHandler,
		AppVersion:     c.AppVersion,
		BrokerProducer: c.BrokerProducer,
		Enforcer:       c.Enforcer,
	}
}
