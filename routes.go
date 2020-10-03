package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dgraph-io/dgo/protos/api"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func Route() *chi.Mux {
	mux := chi.NewMux()
	//Globals middlewares
	mux.Use(
		middleware.Logger,
		middleware.Recoverer,
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

	var buyers Buyer
	if err := json.Unmarshal(resp.GetJson(), &buyers); err != nil {
		log.Fatal(err)
	}

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
