package wire

import (
	"go-api/handler"
	"go-api/infrastructure/config"

	"gorm.io/gorm"
)

type Container struct {
	MetadataHandler *handler.MetadataHandler
}

func NewContainer(db *gorm.DB, env *config.Config) *Container {
	metadataHandler := handler.NewMetadataHandler()

	return &Container{
		MetadataHandler: metadataHandler,
	}
}
