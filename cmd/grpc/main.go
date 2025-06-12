package main

import (
	"log"
	"net"

	"github.com/ruslanukhlin/SwiftTalk.auth-service/pkg/config"
	"github.com/ruslanukhlin/SwiftTalk.auth-service/pkg/gorm"
	pb "github.com/ruslanukhlin/SwiftTalk.common/gen/auth"
	"google.golang.org/grpc"

	authApp "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/application/auth"
	"github.com/ruslanukhlin/SwiftTalk.auth-service/internal/infrastructure/db/postgres"
	authGRPC "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/infrastructure/grpc"
	jwtRepo "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/infrastructure/jwt"
	passwordRepo "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/infrastructure/password"
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
	passwordRepo := passwordRepo.NewPasswordRepo()

	authApp := authApp.NewAuthApp(userRepo, passwordRepo, tokenRepo)

	runGRPCServer(authApp)
}

func runGRPCServer(authApp authApp.AuthService) {
	cfg := config.LoadConfigFromEnv()

	lis, err := net.Listen("tcp", ":" + cfg.PortGrpc)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}

	userGRPCHandler := authGRPC.NewUserGRPCHandler(authApp)
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, userGRPCHandler)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка grpc сервера: %v", err)
	}
}