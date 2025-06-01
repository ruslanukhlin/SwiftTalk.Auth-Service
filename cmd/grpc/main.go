package main

import (
	"log"
	"net"

	application "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/application/user"
	"github.com/ruslanukhlin/SwiftTalk.auth-service/internal/infrastructure/db/postgres"
	"github.com/ruslanukhlin/SwiftTalk.auth-service/pkg/config"
	"github.com/ruslanukhlin/SwiftTalk.auth-service/pkg/gorm"
	"google.golang.org/grpc"

	userGRPC "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/infrastructure/db/grpc"
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

	userDb := postgres.NewPostgresMemoryRepository(gorm.DB)
	userApp := application.NewUserApp(userDb)

	runGRPCServer(userApp, cfg.Port)
}

func runGRPCServer(userApp *application.UserApp, port string) {
	lis, err := net.Listen("tcp", ":" + port)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}

	userGRPCHandler := userGRPC.NewUserGRPCHandler(userApp)
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, userGRPCHandler)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка grpc сервера: %v", err)
	}
}