package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/File-Sharer/user-service/hasher_pbs"
	"github.com/File-Sharer/user-service/internal/config"
	"github.com/File-Sharer/user-service/internal/handler"
	"github.com/File-Sharer/user-service/internal/rabbitmq"
	"github.com/File-Sharer/user-service/internal/repository"
	"github.com/File-Sharer/user-service/internal/repository/postgres"
	"github.com/File-Sharer/user-service/internal/server"
	"github.com/File-Sharer/user-service/internal/service"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing config: %s", err.Error())
	}

	if err := initEnv(); err != nil {
		log.Fatalf("error initializing env: %s", err.Error())
	}

	hasherConn, err := grpc.NewClient(viper.GetString("hasherService.target"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error connecting to hasher grpc service: %s", err.Error())
	}
	defer func ()  {
		if err := hasherConn.Close(); err != nil {
			log.Fatalf("error closing grpc hasher service connection: %s", err.Error())
		}
	}()

	hasherClient := pb.NewHasherClient(hasherConn)

	dbConfig := &config.DBConfig{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		DBName: os.Getenv("DB_NAME"),
		SSLMode: os.Getenv("DB_SSLMODE"),
	}
	db, err := postgres.NewPostgresDB(context.Background(), dbConfig)
	if err != nil {
		log.Fatalf("error opening db: %s", err.Error())
	}
	defer func ()  {
		if err := db.Close(context.Background()); err != nil {
			log.Fatalf("error closing db connection: %s", err.Error())
		}
	}()

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	defer func ()  {
		if err := rdb.Close(); err != nil {
			log.Fatalf("error closing redis connection: %s", err.Error())
		}
	}()

	rabbitmq, err := rabbitmq.New(os.Getenv("RABBITMQ_URI"))
	if err != nil {
		log.Fatalf("error connecting to rabbitmq: %s", err.Error())
	}

	repo := repository.New(db, rdb)
	services := service.New(repo, rabbitmq, hasherClient)
	handlers := handler.New(services, hasherClient)

	srv := server.New()
	serverConfig := &config.ServerConfig{
		Port: viper.GetString("app.port"),
		Handler: handlers.InitRoutes(),
		MaxHeaderBytes: 1 << 20,
		ReadTimeout: time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	go func ()  {
		if err := srv.Run(serverConfig); err != nil {
			log.Fatalf("error occured while running server: %s", err.Error())
		}
	}()

	log.Println("User Server Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("User Server Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("error shutting down server: %s", err.Error())
	}
}

func initConfig() error {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func initEnv() error {
	return godotenv.Load(".env")
}
