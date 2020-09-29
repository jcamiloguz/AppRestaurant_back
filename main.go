package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

//StrucingData infos
type StrucingData struct {
	transactionID string
	Buyer         Buyer
	IP            string
	Device        string
	Product       []Product
}

func main() {
	data := GetRestaInfo("213213")
	fmt.Println(data)

}

func setupAPI() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Get("/buyers", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("home"))
	})

	http.ListenAndServe(":3000", r)
}
