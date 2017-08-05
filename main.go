package main

import (
	"context"
	"github.com/baopham/gominder/api"
	"github.com/baopham/gominder/models"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"os"
)

func main() {
	session, err := models.NewSession()

	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	r.Path("/reminders").
		Methods("GET").
		Handler(Adapt(http.HandlerFunc(api.GetReminders), withDB(session)))

	r.Path("/reminders").
		Methods("POST").
		Handler(Adapt(http.HandlerFunc(api.AddReminder), withDB(session)))

	r.Path("/reminders/{id}").
		Methods("GET").
		Handler(Adapt(http.HandlerFunc(api.GetReminder), withDB(session)))

	r.Path("/remind").
		Methods("PUT").
		Handler(Adapt(http.HandlerFunc(api.Remind), withDB(session)))

	port := "3000"

	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	handler := cors.Default().Handler(r)

	log.Fatalln(http.ListenAndServe("localhost:"+port, handler))
}

// Adapter pattern: https://medium.com/@matryer/production-ready-mongodb-in-go-for-beginners-ef6717a77219
type Adapter func(http.Handler) http.Handler

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func withDB(db *mgo.Session) Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := db.Copy()
			defer session.Close()

			ctx := context.WithValue(r.Context(), "database", session)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
