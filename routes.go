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

const queryBuyers string = `
{
	buyers(func:has(name)){
		 id
		 name
		 age
	}
}
`
const queryDetails string = `
{
	buyer(func: eq(id,"%s")){
		ID as id
		name 
		age
	}
	transaction(func: eq(buyer_id,val(ID))){
		IP as ip
		product_id
		device
		products_id
	}
	var(func: eq (ip, val(IP))){
		BUY as buyer_id
	}
		buyers(func: eq(id, val(BUY))){
			id
			name
			age
			
		}
}
`
const queryProducts string = `
{
	products(func: eq(product_id,"%s")){
		product_id
		product_name
		price
	}

}
`

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

	mux.Post("/sync", postBuyersHandler)
	mux.Route("/buyer", func(r chi.Router) {
		r.With(paginate).Get("/", getBuyersHandler)

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

	resp, err := txn.Query(context.Background(), queryBuyers)
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
		next.ServeHTTP(w, r)
	})
}

func GetBuyer(w http.ResponseWriter, r *http.Request) {
	buyer := chi.URLParam(r, "buyer")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("done-by", "jcamiloguz")
	dgClient := newClient()
	txnBuy := dgClient.NewTxn()
	txnPro := dgClient.NewTxn()
	q := fmt.Sprintf(queryDetails, buyer)
	resp, err := txnBuy.Query(context.Background(), q)
	if err != nil {
		log.Fatal(err)
	}
	var Details Details
	var history History
	if err := json.Unmarshal(resp.GetJson(), &Details); err != nil {
		log.Fatal(err)
	}
	for _, trans := range Details.Transaction {
		for _, product := range trans.Products_id {
			var decodeProduct History
			p := fmt.Sprintf(queryProducts, product)
			resp, err := txnPro.Query(context.Background(), p)
			if err != nil {
				log.Fatal(err)
			}
			if err := json.Unmarshal(resp.GetJson(), &decodeProduct); err != nil {
				log.Fatal(err)
			}
			product := decodeProduct.Products
			history.Products = append(history.Products, product...)

		}
	}
	Response := DetailResponse{
		Details,
		history,
	}
	json.NewEncoder(w).Encode(Response)

}
