package wire

import (
	"go-api/handler"

	"gorm.io/gorm"
)

type Container struct {
	HeuristicHandler *handler.HeuristicHandler
}

func NewContainer(_ *gorm.DB) *Container {
	return &Container{
		HeuristicHandler: handler.NewHeuristicHandler(),
	}
}
