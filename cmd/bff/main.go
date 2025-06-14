package main

// @title SwiftTalk Auth Service API
// @version 1.0
// @description API сервиса аутентификации для платформы SwiftTalk
// @host localhost:8080
// @BasePath /authService/
import (
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/ruslanukhlin/SwiftTalk.auth-service/docs"
	pb "github.com/ruslanukhlin/SwiftTalk.common/gen/auth"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ruslanukhlin/SwiftTalk.auth-service/internal/infrastructure/bff"
	jwtRepo "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/infrastructure/jwt"
	"github.com/ruslanukhlin/SwiftTalk.auth-service/pkg/config"
)

func main() {
	cfg := config.LoadConfigFromEnv()

	app := fiber.New()

	// Swagger route
	app.Get("/swagger/*", fiberSwagger.FiberWrapHandler())

	// Инициализация RSA ключей для JWT
	privateKey, publicKey, err := config.ParseKeys()
	if err != nil {
		log.Fatalf("Ошибка загрузки RSA ключей: %v", err)
	}

	// Инициализация JWT репозитория
	tokenRepo := jwtRepo.NewJWTTokenRepository(privateKey, publicKey)

	conn, err := grpc.NewClient(":" + cfg.PortGrpc, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка подключения к gRPC серверу: %v", err)
	}
	defer conn.Close()

	authClient := pb.NewAuthServiceClient(conn)
	authService := bff.NewAuthService(authClient)
	handler := bff.NewHandler(authService, tokenRepo)

	bff.RegisterRoutes(app, handler)

	if err := app.Listen(":" + cfg.PortHttp); err != nil {
		log.Fatalf("Ошибка запуска HTTP сервера: %v", err)
	}
}