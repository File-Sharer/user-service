package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/File-Sharer/user-service/hasher_pbs"
	"github.com/File-Sharer/user-service/internal/config"
	"github.com/File-Sharer/user-service/internal/handler"
	"github.com/File-Sharer/user-service/internal/repository"
	"github.com/File-Sharer/user-service/internal/repository/postgres"
	"github.com/File-Sharer/user-service/internal/server"
	"github.com/File-Sharer/user-service/internal/service"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}

	if err := initEnv(); err != nil {
		logrus.Fatalf("error initializing env: %s", err.Error())
	}

	hasherConn, err := grpc.NewClient(viper.GetString("hasherService.target"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("error connecting to hasher grpc service: %s", err.Error())
	}
	defer func ()  {
		if err := hasherConn.Close(); err != nil {
			logrus.Fatalf("error closing grpc hasher service connection: %s", err.Error())
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
		logrus.Fatalf("error opening db: %s", err.Error())
	}
	defer func ()  {
		if err := db.Close(context.Background()); err != nil {
			logrus.Fatalf("error closing db connection: %s", err.Error())
		}
	}()

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	defer func ()  {
		if err := rdb.Close(); err != nil {
			logrus.Fatalf("error closing redis connection: %s", err.Error())
		}
	}()

	repo := repository.New(db, rdb)
	services := service.New(repo, hasherClient)
	handlers := handler.New(services)

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
			logrus.Fatalf("error occured while running server: %s", err.Error())
		}
	}()

	logrus.Println("User Server Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Println("User Server Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Fatalf("error shutting down server: %s", err.Error())
	}
}

func initConfig() error {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func initEnv() error {
	return godotenv.Load()
}
