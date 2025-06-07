package main

import (
	"log"
	"net"

	tokenApp "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/application/token"
	userApp "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/application/user"
	cryptRepo "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/infrastructure/crypt"
	"github.com/ruslanukhlin/SwiftTalk.auth-service/internal/infrastructure/db/postgres"
	jwtRepo "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/infrastructure/jwt"
	"github.com/ruslanukhlin/SwiftTalk.auth-service/pkg/config"
	"github.com/ruslanukhlin/SwiftTalk.auth-service/pkg/gorm"
	"google.golang.org/grpc"

	userGRPC "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/infrastructure/grpc"
	pb "github.com/ruslanukhlin/SwiftTalk.common/gen/auth"
)

func main() {
	cfg := config.LoadConfigFromEnv()	

	if err := gorm.InitDB(config.DNS(cfg.Postgres)); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}

	if err := gorm.Migrate(cfg); err != nil {
		log.Fatalf("Ошибка миграции базы данных: %v", err)
	}

	userRepo := postgres.NewPostgresMemoryRepository(gorm.DB)
	tokenRepo := jwtRepo.NewJWTTokenRepository(cfg.JWT)
	passwordRepo := cryptRepo.NewCryptRepository()

	tokenApp := tokenApp.NewTokenApp(tokenRepo)
	userApp := userApp.NewUserApp(userRepo, tokenRepo, passwordRepo)

	runGRPCServer(userApp, tokenApp)
}

func runGRPCServer(userApp *userApp.UserApp, tokenApp *tokenApp.TokenApp) {
	cfg := config.LoadConfigFromEnv()

	lis, err := net.Listen("tcp", ":" + cfg.Port)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}

	userGRPCHandler := userGRPC.NewUserGRPCHandler(userApp, tokenApp)
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, userGRPCHandler)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка grpc сервера: %v", err)
	}
}