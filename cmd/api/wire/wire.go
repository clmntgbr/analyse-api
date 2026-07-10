package wire

import (
	"go-api/handler"
	"go-api/handler/middleware"
	infraClerk "go-api/infrastructure/clerk"
	"go-api/infrastructure/config"
	"go-api/infrastructure/storage"
	repoGorm "go-api/repository/gorm"
	"go-api/usecase/auth"
	"go-api/usecase/clerk"
	"go-api/usecase/media"
	"go-api/usecase/user"
	"log"

	"gorm.io/gorm"
)

type Container struct {
	AuthenticateMiddleware *middleware.AuthenticateMiddleware
	ClerkMiddleware        *middleware.ClerkMiddleware
	ClerkHandler           *handler.ClerkHandler
	UserHandler            *handler.UserHandler
	MediaHandler           *handler.MediaHandler
}

func NewContainer(db *gorm.DB, env *config.Config) *Container {
	jwksProvider, err := infraClerk.NewJWKSProvider(env)
	if err != nil {
		log.Fatalf("failed to create JWKS provider: %v", err)
	}
	log.Println("🚀 JWKS provider created")

	storageClient, err := storage.NewMinIOStorage(env)
	if err != nil {
		log.Fatalf("failed to create storage client: %v", err)
	}
	log.Println("🚀 Storage client created")

	generatePresignedUploadUrlUseCase := media.NewGeneratePresignedUploadUrlUseCase(storageClient)

	userRepo := repoGorm.NewUserRepository(db)

	validateTokenUseCase := auth.NewValidateTokenUseCase(jwksProvider, &userRepo)
	fetchUserUseCase := clerk.NewFetchUserUseCase(env)
	getUserByClerkIDUseCase := user.NewGetUserByClerkIDUseCase(&userRepo)
	createUserUseCase := user.NewCreateUserUseCase(&userRepo)
	updateUserUseCase := user.NewUpdateUserUseCase(&userRepo)
	deleteUserByClerkIDUseCase := user.NewDeleteUserByClerkIDUseCase(&userRepo)

	clerkMiddleware := middleware.NewClerkMiddleware(env.ClerkWebhookSecret)
	authenticateMiddleware := middleware.NewAuthenticateMiddleware(
		validateTokenUseCase,
		fetchUserUseCase,
		createUserUseCase,
		updateUserUseCase,
	)

	return &Container{
		AuthenticateMiddleware: authenticateMiddleware,
		ClerkMiddleware:        clerkMiddleware,
		ClerkHandler: handler.NewClerkHandler(
			getUserByClerkIDUseCase,
			createUserUseCase,
			updateUserUseCase,
			deleteUserByClerkIDUseCase,
		),
		UserHandler: handler.NewUserHandler(),
		MediaHandler: handler.NewMediaHandler(
			generatePresignedUploadUrlUseCase,
		),
	}
}
