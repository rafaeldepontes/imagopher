package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/rafaeldepontes/imagopher/internal/handler"
	"github.com/rafaeldepontes/imagopher/internal/tool"
)

func init() {
	env := ".env"
	tool.ChecksEnv(&env)
	if err := godotenv.Load(env); err != nil {
		log.Fatalln("[ERROR] Couldn't load the env variable:", err)
	}
}

func main() {
	port := os.Getenv("PORT")

	r := chi.NewRouter()
	handler.Routes(r)

	log.Printf("Application running on %s\n", "localhost:"+port)
	log.Fatalln(http.ListenAndServe(":"+port, r))
}
