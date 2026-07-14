package wire

import (
	"go-api/handler"
	"go-api/infrastructure/config"
	heuristicsinfra "go-api/infrastructure/heuristics"
	"go-api/infrastructure/messaging/rabbitmq"
	"go-api/infrastructure/messaging/security"
	"go-api/infrastructure/storage"
	repoGorm "go-api/repository/gorm"
	heuristicuc "go-api/usecase/heuristic"
	pipelineuc "go-api/usecase/pipeline"
	"go-api/usecase/signal"
	"log"

	"gorm.io/gorm"
)

type Container struct {
	HeuristicHandler *handler.HeuristicHandler
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

	analyzer := heuristicsinfra.NewAnalyzer()
	analyzeMediaHeuristicsUseCase := heuristicuc.NewAnalyzeMediaHeuristicsUseCase(storageClient, analyzer)
	createSignalUseCase := signal.NewCreateSignalUseCase(&signalRepo)

	parser := security.NewWorkerParser(env)
	securityValidator := security.NewWorkerSecurityValidator(env)

	heuristicHandler := handler.NewHeuristicHandler(
		parser,
		securityValidator,
		dispatcher,
		analyzeMediaHeuristicsUseCase,
		createSignalUseCase,
	)

	return &Container{
		HeuristicHandler: heuristicHandler,
	}
}
