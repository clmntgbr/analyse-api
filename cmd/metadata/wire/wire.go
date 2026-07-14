package wire

import (
	"go-api/handler"
	"go-api/infrastructure/config"
	"go-api/infrastructure/centrifugo"
	"go-api/infrastructure/messaging/rabbitmq"
	"go-api/infrastructure/messaging/security"
	metadatainfra "go-api/infrastructure/metadata"
	"go-api/infrastructure/storage"
	repoGorm "go-api/repository/gorm"
	metadatauc "go-api/usecase/metadata"
	pipelineuc "go-api/usecase/pipeline"
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

	publisher, err := rabbitmq.NewPublisherFromEnv(env)
	if err != nil {
		log.Fatalf("failed to create publisher: %v", err)
	}

	centrifugoPublisher := centrifugo.NewPublisher(env)

	mediaRepo := repoGorm.NewMediaRepository(db)
	signalRepo := repoGorm.NewSignalRepository(db)

	aggregateAnalysisUseCase := pipelineuc.NewAggregateAnalysisUseCase(&mediaRepo, &signalRepo, centrifugoPublisher)
	dispatcher := pipelineuc.NewDispatcher(env, &mediaRepo, publisher, aggregateAnalysisUseCase)

	analyzer := metadatainfra.NewAnalyzer()
	analyzeMediaMetadataUseCase := metadatauc.NewAnalyzeMediaMetadataUseCase(storageClient, analyzer)
	createSignalUseCase := signal.NewCreateSignalUseCase(&signalRepo)

	parser := security.NewWorkerParser(env)
	securityValidator := security.NewWorkerSecurityValidator(env)

	metadataHandler := handler.NewMetadataHandler(
		parser,
		securityValidator,
		dispatcher,
		analyzeMediaMetadataUseCase,
		createSignalUseCase,
	)

	return &Container{
		MetadataHandler: metadataHandler,
	}
}
