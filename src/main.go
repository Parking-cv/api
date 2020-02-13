package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
	"parking-cv/server/src/valet"
	"time"
)

func main() {

	err := valet.InitializeMongoClient()
	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Mount("/api", api())
	router.Mount("/pi", piRouter())

	_ = http.ListenAndServe(":4321", router)
}

// Router for requests from picameras
func piRouter() http.Handler {
	router := chi.NewRouter()
	router.Use(PiAuth)

	router.Post("/frames", receiveFrames)

	return router
}

// Router for interface rest api
func api() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Timeout(60 * time.Second))
	router.Get("/test", testRoute)

	return router
}

func PiAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify Pi's JWT
		next.ServeHTTP(w, r)
	})
}
