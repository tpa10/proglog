package server

import (
    "encoding/json"
    "net/http"

    "github.com/gorilla/mux"
)

// Used to create an instance of an HTTP server
func NewHTTPSServer(addr string) *http.Server {
    httpsrv := newHTTPServer()
    r := mux.NewRouter()
    r.HandleFunc("/", httpsrv.handleProduce).Methods("POST")
    r.HandleFunc("/", httpsrv.handleConsume).Methods("GET")
    return &http.Server{
        Addr:   addr,
        Handler: r,
    }     
}

type <pick up here - tpa>
