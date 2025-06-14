package main

// @title SwiftTalk Auth Service API
// @version 1.0
// @description API сервиса аутентификации для платформы SwiftTalk
// @host localhost:8080
// @BasePath /authService/
import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	_ "github.com/ruslanukhlin/SwiftTalk.auth-service/docs"
	pb "github.com/ruslanukhlin/SwiftTalk.common/gen/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ruslanukhlin/SwiftTalk.auth-service/internal/infrastructure/bff"
	jwtRepo "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/infrastructure/jwt"
	"github.com/ruslanukhlin/SwiftTalk.auth-service/pkg/config"
)

func main() {
	cfg := config.LoadConfigFromEnv()

	conn, err := grpc.NewClient(":"+cfg.PortGrpc, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка подключения к GRPC серверу: %v", err)
	}

	// Перемещаем defer после всех инициализаций
	server := fiber.New()
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
