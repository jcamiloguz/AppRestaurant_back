package main

import (
	"encoding/json"
	"net/http"

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
	return mux
}
func getBuyersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("done-by", "jcamiloguz")
	info := GetRestaInfo("2332132")
	json.NewEncoder(w).Encode(info)

}
