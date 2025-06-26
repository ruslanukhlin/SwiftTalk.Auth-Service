package main

// @title SwiftTalk Auth Service API
// @version 1.0
// @description API сервиса аутентификации для платформы SwiftTalk
// @host localhost:5002
// @BasePath /auth-service/
import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/ruslanukhlin/SwiftTalk.Auth-service/docs"
	pb "github.com/ruslanukhlin/SwiftTalk.Common/gen/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ruslanukhlin/SwiftTalk.Auth-service/internal/infrastructure/bff"
	jwtRepo "github.com/ruslanukhlin/SwiftTalk.Auth-service/internal/infrastructure/jwt"
	"github.com/ruslanukhlin/SwiftTalk.Auth-service/pkg/config"
)

func main() {
	cfg := config.LoadConfigFromEnv()

	conn, err := grpc.NewClient(":"+cfg.PortGrpc, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка подключения к GRPC серверу: %v", err)
	}

	server := fiber.New()

	server.Get("/docs/docs.json", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger-v3.json")
	})

	// OpenAPI 3.0 UI
	server.Get("/docs/*", swagger.New(swagger.Config{
		URL:          "docs.json",
		DeepLinking:  true,
		DocExpansion: "none",
		Title:        "SwiftTalk Auth Service API (OpenAPI 3.0)",
	}))

	authClient := pb.NewAuthServiceClient(conn)
	privateKey, publicKey, err := config.ParseKeys()
	if err != nil {
		log.Fatalf("Ошибка загрузки RSA ключей: %v", err)
	}
	tokenRepo := jwtRepo.NewJWTTokenRepository(privateKey, publicKey)
	authService := bff.NewAuthService(authClient)
	handler := bff.NewHandler(authService, tokenRepo)
	bff.RegisterRoutes(server, handler)

	go func() {
		if err := server.Listen(":" + cfg.PortHttp); err != nil {
			log.Printf("Ошибка запуска HTTP сервера: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	if err := conn.Close(); err != nil {
		log.Printf("Ошибка закрытия GRPC соединения: %v", err)
	}
	if err := server.Shutdown(); err != nil {
		log.Printf("Ошибка graceful shutdown: %v", err)
	}
}
