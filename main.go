package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/rafaeldepontes/imagopher/internal/application"
	"github.com/rafaeldepontes/imagopher/internal/handler"
	"github.com/rafaeldepontes/imagopher/internal/tool"
)

var app application.Application

func init() {
	env := ".env"
	tool.ChecksEnv(&env)
	if err := godotenv.Load(env); err != nil {
		log.Fatalln("[ERROR] Couldn't load the env variable:", err)
	}

	app = application.NewApplication()
}

// TODO: Let's be real here... who needs an interface when you have the source code
// and the cli using curl (or Postman... I like postman.).
//
// I'm changed my mind. I'll do something close to pinterest. (But still... None login needed).
func main() {
	port := os.Getenv("PORT")

	r := chi.NewRouter()
	handler.Routes(r, app)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("[INFO] Application running on %s\n", "localhost:"+port)
		log.Fatalln("[ERROR] ", http.ListenAndServe(":"+port, r))
	}()

	log.Println("[INFO] Received shutdown signal. Shutting down...")

	<-sigChan

	if err := app.Shutdown(); err != nil {
		log.Fatalf("[ERROR] Could not shutdown: %v\n", err)
	}

	log.Println("[INFO] Shutdown complete.")
}
