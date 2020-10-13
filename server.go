package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

//MyServer own struct server
type MyServer struct {
	server *http.Server
}

//NewServer Create a CHI API Server
func NewServer(mux *chi.Mux) *MyServer {
	s := &http.Server{
		Addr:           ":9000",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return &MyServer{s}
}

//Run a Server
func (s *MyServer) Run() {
	log.Fatal(s.server.ListenAndServe())

}
