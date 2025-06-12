package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	pb "github.com/ruslanukhlin/SwiftTalk.common/gen/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ruslanukhlin/SwiftTalk.auth-service/internal/infrastructure/bff"
	"github.com/ruslanukhlin/SwiftTalk.auth-service/pkg/config"
)

func main() {
	cfg := config.LoadConfigFromEnv()

	app := fiber.New()

	conn, err := grpc.NewClient(":" + cfg.PortGrpc, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка подключения к gRPC серверу: %v", err)
	}
	defer conn.Close()

	authClient := pb.NewAuthServiceClient(conn)
	authService := bff.NewAuthService(authClient)
	handler := bff.NewHandler(authService)

	bff.RegisterRoutes(app, handler)

	if err := app.Listen(":" + cfg.PortHttp); err != nil {
		log.Fatalf("Ошибка запуска HTTP сервера: %v", err)
	}
}