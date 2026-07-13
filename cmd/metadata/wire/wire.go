package wire

import (
	"go-api/handler"
	"go-api/infrastructure/config"
	"go-api/infrastructure/messaging/security"
	metadatainfra "go-api/infrastructure/metadata"
	"go-api/infrastructure/storage"
	repoGorm "go-api/repository/gorm"
	metadatauc "go-api/usecase/metadata"
	"go-api/usecase/signal"
	"log"

	"gorm.io/gorm"
)

type Container struct {
	MetadataHandler *handler.MetadataHandler
}

func NewContainer(db *gorm.DB, env *config.Config) *Container {
	storageClient, err := storage.NewMinIOStorage(env)
	if err != nil {
		log.Fatalf("failed to create storage client: %v", err)
	}

	signalRepo := repoGorm.NewSignalRepository(db)

	analyzer := metadatainfra.NewAnalyzer()
	analyzeMediaMetadataUseCase := metadatauc.NewAnalyzeMediaMetadataUseCase(storageClient, analyzer)
	createSignalUseCase := signal.NewCreateSignalUseCase(&signalRepo)

	metadataHandler := handler.NewMetadataHandler(
		env,
		security.NewWorkerParser(env),
		security.NewWorkerSecurityValidator(env),
		analyzeMediaMetadataUseCase,
		createSignalUseCase,
	)

	return &Container{
		MetadataHandler: metadataHandler,
	}
}
