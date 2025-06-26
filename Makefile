all: run

run:
	go run cmd/grpc/main.go & go run cmd/bff/main.go

swagger:
	swag init -g cmd/bff/main.go -d . -o ./docs

swagger-v3:
	swag init -g cmd/bff/main.go -d . -o ./docs
	swagger2openapi ./docs/swagger.json -o ./docs/swagger-v3.json

docker-up:
	docker compose --env-file .env.local up -d

docker-up prod:
	docker compose --env-file .env.prod up -d