package app

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"finstat/internal/handlers"
	"finstat/internal/repository/balance"
	"finstat/internal/repository/user"
	"finstat/internal/service"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found")
	}
}

func Run() {
	dbuser := os.Getenv("DB_USER")
	dbpass := os.Getenv("DB_PASSWORD")
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	postgresData := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", dbuser, dbpass, dbhost, dbport, dbname)
	postgresDB, err := sqlx.Connect("postgres", postgresData)
	if err != nil {
		log.Fatalln(err)
	}

	balancerepo := balance.NewRepo(postgresDB)
	userrepo := user.NewRepo(postgresDB)

	services := service.NewService(userrepo, balancerepo)

	handler := handlers.NewHandler(services)

	r := handler.Init()

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15,
		"the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m",
	)
	flag.Parse()

	srv := &http.Server{
		Addr:         ":8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	go func() {
		fs := http.FileServer(http.Dir("/app/docs"))
		if err := http.ListenAndServe(":1349", fs); err != nil {
			log.Println(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		err := postgresDB.Close()
		if err != nil {
			log.Printf("database is connected incorrectly: %v\n", err)
		} else {
			log.Println("database is connected correctly")
		}

		cancel()
	}()

	err = srv.Shutdown(ctx)
	if err != nil {
		log.Printf("server is shutdown: %s", err)
	}

	select {
	case <-ctx.Done():
		log.Println("timeout 0f 5 second left")
	}

	log.Println("server is successfully stopped")

}
