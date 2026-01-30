package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/rafaeldepontes/imagopher/internal/application"
	"github.com/rafaeldepontes/imagopher/internal/application/model"
	"github.com/rafaeldepontes/imagopher/internal/handler"
	"github.com/rafaeldepontes/imagopher/internal/tool"
)

var app *model.Application

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
func main() {
	port := os.Getenv("PORT")

	r := chi.NewRouter()
	handler.Routes(r, app)

	log.Printf("Application running on %s\n", "localhost:"+port)
	log.Fatalln(http.ListenAndServe(":"+port, r))
}
