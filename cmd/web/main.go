package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// DB
	"github.com/RehanAthallahAzhar/shopeezy-orders/db"
	"github.com/RehanAthallahAzhar/shopeezy-orders/internal/configs"
	dbGenerated "github.com/RehanAthallahAzhar/shopeezy-orders/internal/db"
	"github.com/RehanAthallahAzhar/shopeezy-orders/internal/handlers"

	// Migrations
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	// Internal Packages
	"github.com/RehanAthallahAzhar/shopeezy-orders/internal/delivery/http/middlewares"
	"github.com/RehanAthallahAzhar/shopeezy-orders/internal/delivery/http/routes"
	"github.com/RehanAthallahAzhar/shopeezy-orders/internal/gateways/messaging"
	"github.com/RehanAthallahAzhar/shopeezy-orders/internal/models"
	"github.com/RehanAthallahAzhar/shopeezy-orders/internal/pkg/grpc/account"
	"github.com/RehanAthallahAzhar/shopeezy-orders/internal/pkg/logger"
	"github.com/RehanAthallahAzhar/shopeezy-orders/internal/pkg/redis"
	"github.com/RehanAthallahAzhar/shopeezy-orders/internal/repositories"
	"github.com/RehanAthallahAzhar/shopeezy-orders/internal/services"

	// Protobuf Packages
	accountpb "github.com/RehanAthallahAzhar/shopeezy-protos/pb/account"
	authpb "github.com/RehanAthallahAzhar/shopeezy-protos/pb/auth"
	productpb "github.com/RehanAthallahAzhar/shopeezy-protos/pb/product"
)

func main() {
	log := logger.NewLogger()
	log.Info("Memulai Order Service...")

	cfg, err := configs.LoadConfig(log)
	if err != nil {
		log.Fatalf("FATAL: Gagal memuat konfigurasi: %v", err)
	}

	log.Printf(">>>> BUKTI ALAMAT GRPC YANG AKAN DIGUNAKAN: '%s'", cfg.GRPC.AccountServiceAddress)

	dbCredential := models.Credential{
		Host:         cfg.Postgre.Host,
		Username:     cfg.Postgre.User,
		Password:     cfg.Postgre.Password,
		DatabaseName: cfg.Postgre.Name,
		Port:         cfg.Postgre.Port,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := db.Connect(ctx, &dbCredential)
	if err != nil {
		log.Errorf("DB connection error: %v", err)
	}

	// Migration
	log.Println("Running database migrations...")
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbCredential.Username,
		dbCredential.Password,
		dbCredential.Host,
		dbCredential.Port,
		dbCredential.DatabaseName,
	)
	m, err := migrate.New(
		cfg.Migration.Path,
		connectionString,
	)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to execute database migrations: %v", err)
	}

	log.Println("Database migration successfully executed.")

	// Init SQLC Store untuk Transaksi
	sqlcQueries := dbGenerated.New(conn)
	store := dbGenerated.NewStore(conn)

	// Redis
	redisClient, err := redis.NewRedisClient()
	if err != nil {
		log.Fatalf("Failed to create Redis client: %v", err)
	}
	defer redisClient.Close()

	// gRPC Product
	productConn := createGrpcConnection(cfg.GRPC.ProductServiceAddress, log)
	defer productConn.Close()

	productClient := productpb.NewProductServiceClient(productConn)
	log.Printf("Product Service Connected to %s", cfg.GRPC.ProductServiceAddress)

	// gRPC Account & Auth
	accountConn := createGrpcConnection(cfg.GRPC.AccountServiceAddress, log)
	defer accountConn.Close()
	accountClient := accountpb.NewAccountServiceClient(accountConn)

	authClient := authpb.NewAuthServiceClient(accountConn)
	authClientWrapper := account.NewAuthClientFromService(authClient, accountConn)

	log.Printf("Account Service Connected to %s", cfg.GRPC.AccountServiceAddress)

	// Publisher Rabbitmq
	rabbitMQURL := cfg.RabbitMQ.URL
	rabbitConn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()

	rabbitChannel, err := rabbitConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open RabbitMQ channel: %v", err)
	}
	defer rabbitChannel.Close()

	orderPublisher, err := messaging.NewRabbitMQPublisher(rabbitChannel)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ publisher: %v", err)
	}

	// Dependency Injection
	validate := validator.New()
	orderRepo := repositories.NewOrderRepository(conn, sqlcQueries, store, log)
	orderService := services.NewOrderService(orderRepo, redisClient, productClient, accountClient, orderPublisher, validate, log)
	orderHandler := handlers.NewHandler(orderService, log)

	// midlleware
	authMiddleware := middlewares.AuthMiddleware(authClientWrapper, log)

	// Setup Server Web (Echo)
	e := echo.New()

	routes.InitRoutes(e, orderHandler, authMiddleware)

	log.Printf("Server Running on port%s", cfg.Server.Port)
	if err := e.Start(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func createGrpcConnection(url string, log *logrus.Logger) *grpc.ClientConn {
	// Gunakan grpc.Dial, yang modern dan non-blocking
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		// Gunakan Fatalf di sini karena jika koneksi gagal saat startup,
		// aplikasi tidak bisa berjalan dengan benar.
		log.Fatalf("Failed to connect to gRPC service at %s: %v", url, err)
	}
	log.Printf("Connected to gRPC service at %s", url)
	return conn
}
