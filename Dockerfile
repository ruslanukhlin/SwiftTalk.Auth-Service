FROM golang:1.24.4-alpine AS build

WORKDIR /build

COPY . .

RUN go mod download && go mod verify

WORKDIR /build/cmd/grpc

RUN go build -o auth-service-grpc .

WORKDIR /build/cmd/bff

RUN go build -o auth-service-bff .

FROM alpine

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=build /build/.env.prod /app/.env

# Нужно сделать так чтобы она генерировалась в docker file, а не копировалась из локальной машины
COPY --from=build /build/docs /app/docs
COPY --from=build /build/config /app/config
COPY --from=build /build/cmd/grpc/auth-service-grpc /app/auth-service-grpc
COPY --from=build /build/cmd/bff/auth-service-bff /app/auth-service-bff

RUN chmod +x /app/auth-service-grpc /app/auth-service-bff

EXPOSE 50052
EXPOSE 5002

CMD ["/bin/sh", "-c", "/app/auth-service-grpc & /app/auth-service-bff"]