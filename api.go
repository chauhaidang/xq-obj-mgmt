package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	addr  string
	store Store
}

func NewApiServer(addr string, store Store) *APIServer {
	return &APIServer{
		addr: addr, store: store,
	}
}

func (s *APIServer) Serve() {
	r := mux.NewRouter()
	sr := r.PathPrefix("/api/v1").Subrouter()

	objSvc := NewObjectsService(s.store)
	objSvc.RegisterRoutes(sr)

	usrSvc := NewUsersService(s.store)
	usrSvc.RegisterRoutes(sr)

	log.Printf("Starting server at: %v", s.addr)
	log.Fatal(http.ListenAndServe(s.addr, sr))
}
