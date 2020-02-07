package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

//var DB *sql.DB

func main() {

	//db, err := sql.Open("postgres", os.Getenv("POSTGRES_URI"))
	//if err != nil {
	//	panic("Cannot connect to database. Check your connection URI and try again.")
	//}
	//
	//DB = db

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Route("/", func(index chi.Router) {
		index.Post("/test", testRoute)
		index.Post("/frame", analyzeEvent)
	})

	_ = http.ListenAndServe(":3000", router)
}
