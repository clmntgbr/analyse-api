package wire

import (
	"go-api/handler"
	"go-api/infrastructure/aimodel"
	"go-api/infrastructure/config"
	"go-api/infrastructure/messaging/rabbitmq"
	"go-api/infrastructure/messaging/security"
	"go-api/infrastructure/sightengine"
	"go-api/infrastructure/storage"
	repoGorm "go-api/repository/gorm"
	aimodeluc "go-api/usecase/aimodel"
	pipelineuc "go-api/usecase/pipeline"
	"go-api/usecase/signal"
	"log"

	"gorm.io/gorm"
)

type Container struct {
	AiModelHandler *handler.AiModelHandler
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

	mediaRepo := repoGorm.NewMediaRepository(db)
	signalRepo := repoGorm.NewSignalRepository(db)

	aggregateAnalysisUseCase := pipelineuc.NewAggregateAnalysisUseCase(&mediaRepo, &signalRepo)
	dispatcher := pipelineuc.NewDispatcher(env, &mediaRepo, publisher, aggregateAnalysisUseCase)

	sightengineClient := sightengine.NewClient(env)
	analyzer := aimodel.NewAnalyzer(sightengineClient)
	analyzeMediaAiModelUseCase := aimodeluc.NewAnalyzeMediaAiModelUseCase(storageClient, analyzer)
	createSignalUseCase := signal.NewCreateSignalUseCase(&signalRepo)

	parser := security.NewWorkerParser(env)
	securityValidator := security.NewWorkerSecurityValidator(env)

	aiModelHandler := handler.NewAiModelHandler(
		parser,
		securityValidator,
		dispatcher,
		analyzeMediaAiModelUseCase,
		createSignalUseCase,
	)

	return &Container{
		AiModelHandler: aiModelHandler,
	}
}
