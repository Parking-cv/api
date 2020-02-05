package main

import (
	"context"
	"database/sql"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
	"os"
)

var DB *sql.DB

func main() {

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URI"))
	if err != nil {
		panic("Cannot connect to database. Check your connection URI and try again.")
	}

	DB = db

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Route("/", func(index chi.Router) {
		index.Use(dbCtx)
		index.Post("/entry", handleEntry)
		index.Post("/exit", handleExit)
	})

	_ = http.ListenAndServe(":3000", router)
}

func dbCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx := context.WithValue(req.Context(), "DB", DB)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
