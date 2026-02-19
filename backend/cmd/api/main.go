package main

import (
	"backend/internal/migrations"
	"backend/internal/server"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func gracefulShutdown(fiberServer *server.FiberServer, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := fiberServer.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")
	done <- true
}

func main() {
	// Run migrations
	databaseURL := buildDatabaseURL()
	if err := migrations.RunMigrations(databaseURL); err != nil {
		log.Printf("Migration warning: %v", err)
	}

	server := server.New()
	server.RegisterFiberRoutes()

	done := make(chan bool, 1)

	go func() {
		port, _ := strconv.Atoi(os.Getenv("PORT"))
		if port == 0 {
			port = 8080
		}
		err := server.Listen(fmt.Sprintf(":%d", port))
		if err != nil {
			panic(fmt.Sprintf("http server error: %s", err))
		}
	}()

	go gracefulShutdown(server, done)

	<-done
	log.Println("Graceful shutdown complete.")
}

func buildDatabaseURL() string {
	host := getEnv("RISK_REGISTER_DB_HOST", "localhost")
	port := getEnv("RISK_REGISTER_DB_PORT", "5432")
	user := getEnv("RISK_REGISTER_DB_USERNAME", "postgres")
	password := getEnv("RISK_REGISTER_DB_PASSWORD", "postgres")
	database := getEnv("RISK_REGISTER_DB_DATABASE", "risk_register")
	schema := getEnv("RISK_REGISTER_DB_SCHEMA", "public")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		user, password, host, port, database, schema)
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
