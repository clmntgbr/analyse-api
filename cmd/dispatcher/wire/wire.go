package wire

import (
	"go-api/handler"
	"go-api/infrastructure/config"
	"go-api/infrastructure/messaging/rabbitmq"
	"go-api/infrastructure/messaging/security"
	repoGorm "go-api/repository/gorm"
	pipelineuc "go-api/usecase/pipeline"
	"log"

	"gorm.io/gorm"
)

type Container struct {
	AnalyzeRequestHandler *handler.AnalyzeRequestHandler
	StageDoneHandler      *handler.StageDoneHandler
	WorkerPool            *rabbitmq.WorkerPool
}

func NewContainer(db *gorm.DB, env *config.Config) *Container {
	publisher, err := rabbitmq.NewPublisherFromEnv(env)
	if err != nil {
		log.Fatalf("failed to create publisher: %v", err)
	}

	mediaRepo := repoGorm.NewMediaRepository(db)
	signalRepo := repoGorm.NewSignalRepository(db)

	aggregateAnalysisUseCase := pipelineuc.NewAggregateAnalysisUseCase(&mediaRepo, &signalRepo)
	dispatcher := pipelineuc.NewDispatcher(env, &mediaRepo, publisher, aggregateAnalysisUseCase)

	parser := security.NewWorkerParser(env)
	securityValidator := security.NewWorkerSecurityValidator(env)

	analyzeRequestHandler := handler.NewAnalyzeRequestHandler(parser, securityValidator, dispatcher)
	stageDoneHandler := handler.NewStageDoneHandler(parser, securityValidator, dispatcher)

	workerPool := rabbitmq.NewDispatcherWorkers(env, analyzeRequestHandler, stageDoneHandler)

	return &Container{
		AnalyzeRequestHandler: analyzeRequestHandler,
		StageDoneHandler:      stageDoneHandler,
		WorkerPool:            workerPool,
	}
}
