package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

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

	privateKey, publicKey, err := config.ParseKeys()
	if err != nil {
		log.Fatalf("Ошибка загрузки RSA ключей: %v", err)
	}

	userRepo := postgres.NewPostgresMemoryRepository(gorm.DB)
	passwordRepo := passwordRepo.NewPasswordRepo()
	tokenRepo := jwtRepo.NewJWTTokenRepository(privateKey, publicKey)

	authApp := authApp.NewAuthApp(userRepo, passwordRepo, tokenRepo)

	runGRPCServer(authApp)
}

func runGRPCServer(authApp authApp.AuthService) {
	cfg := config.LoadConfigFromEnv()

	lis, err := net.Listen("tcp", ":"+cfg.PortGrpc)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}

	userGRPCHandler := authGRPC.NewUserGRPCHandler(authApp)
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, userGRPCHandler)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("Ошибка grpc сервера: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	grpcServer.GracefulStop()
	log.Println("Сервер успешно остановлен")
}
