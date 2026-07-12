package wire

import (
	"go-api/handler"
	"go-api/infrastructure/config"

	"go-api/infrastructure/messaging/security"

	"gorm.io/gorm"
)

type Container struct {
	MetadataHandler *handler.MetadataHandler
}

func NewContainer(db *gorm.DB, env *config.Config) *Container {
	metadataHandler := handler.NewMetadataHandler(
		env,
		security.NewWorkerParser(env),
		security.NewWorkerSecurityValidator(env),
	)

	return &Container{
		MetadataHandler: metadataHandler,
	}
}
