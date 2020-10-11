package main

import (
	"context"
	"encoding/json"
	"fmt"
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

	mux.Route("/buyer", func(r chi.Router) {
		r.With(paginate).Get("/", getBuyersHandler)
		r.With(paginate).Post("/", postBuyersHandler)

		r.Route("/{buyer}", func(r chi.Router) {
			r.Get("/", GetBuyer) // GET /articles/123
		})

	})

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
	txnB := dgClient.NewTxn()
	txnP := dgClient.NewTxn()
	txnT := dgClient.NewTxn()
	date := r.URL.Query().Get("date")
	buyers, products, transactions := GetRestaInfo(date)
	buyersJSON, productsJSON, transactionsJSON := JSONMarshall(buyers, products, transactions)

	mutBuyers := &api.Mutation{
		SetJson: buyersJSON,
	}
	mutProds := &api.Mutation{
		SetJson: productsJSON,
	}
	mutTrans := &api.Mutation{
		SetJson: transactionsJSON,
	}

	reqBuyers := &api.Request{CommitNow: true, Mutations: []*api.Mutation{mutBuyers}}
	reqProds := &api.Request{CommitNow: true, Mutations: []*api.Mutation{mutProds}}
	reqTrans := &api.Request{CommitNow: true, Mutations: []*api.Mutation{mutTrans}}

	resBuyers, err := txnB.Do(context.Background(), reqBuyers)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(resBuyers)

	resProds, err := txnP.Do(context.Background(), reqProds)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(resProds)

	resTrans, err := txnT.Do(context.Background(), reqTrans)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(resTrans)

}
func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// just a stub.. some ideas are to look at URL query params for something like
		// the page number, or the limit, and send a query cursor down the chain
		next.ServeHTTP(w, r)
	})
}

func GetBuyer(w http.ResponseWriter, r *http.Request) {
	// Assume if we've reach this far, we can access the article
	// context because this handler is a child of the ArticleCtx
	// middleware. The worst case, the recoverer middleware will save us.
	buyer := chi.URLParam(r, "buyer")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("done-by", "jcamiloguz")
	dgClient := newClient()
	txn := dgClient.NewTxn()
	q := fmt.Sprintf(`{
		transaction(func: eq(buyer_id,"%s")){
      ip
      product_id
      device
      products_id
		}
	}`, buyer)
	resp, err := txn.Query(context.Background(), q)
	if err != nil {
		log.Fatal(err)
	}
	var Transactions RespTransaction
	var TransactionsResp TransactionsRsp
	error := json.Unmarshal(resp.GetJson(), &Transactions)
	if error != nil {
		fmt.Println("error:", error)
	}
	data := Transactions.Transactions
	json.NewEncoder(w).Encode(data)

	errorH := json.Unmarshal(data, &TransactionsResp)
	if errorH != nil {
		fmt.Println("error2:", errorH)
	}
	fmt.Printf("%+v", TransactionsResp)

}
