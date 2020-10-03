package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dgraph-io/dgo/protos/api"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func Route() *chi.Mux {
	mux := chi.NewMux()
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	//Globals middlewares
	mux.Use(
		middleware.Logger,
		middleware.Recoverer,
		cors.Handler,
	)
	mux.Get("/buyer", getBuyersHandler)
	mux.Post("/buyer", postBuyersHandler)
	return mux
}
func getBuyersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("done-by", "jcamiloguz")
	dgClient := newClient()
	txn := dgClient.NewTxn()
	const q = `
	{
		buyers(func:has(name)){
			 id
			 name
			 age
		}
	}
	`
	resp, err := txn.Query(context.Background(), q)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(resp.Json)
}

func postBuyersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("done-by", "jcamiloguz")
	dgClient := newClient()
	txn := dgClient.NewTxn()
	date := r.URL.Query().Get("date")
	info := GetRestaInfo(date)
	infojs, err := json.Marshal(info)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(info)
	mu := &api.Mutation{
		SetJson: infojs,
	}

	req := &api.Request{CommitNow: true, Mutations: []*api.Mutation{mu}}

	res, err := txn.Do(context.Background(), req)

	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(res)

}
