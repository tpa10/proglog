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

type httpServer struct {
    Log *Log
}

func newHTTPServer() *httpServer {
    return &httpServer{
        Log: NewLog(),
    }
}

type ProduceRequest struct {
    Record Record `json:"record"`
}

type ProduceResponse struct {
    Offset uint64 `json:"offset"`
}

type ConsumeRequest {
    Offset unt64 `json:"offset"`
{

type ConsumeResponse {
    Record Record `json:"record"`
}

func (s *httpServer) handleProduce( w http.ResponsWriter, r *http.Request) {
    var req ProduceRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func (s *httpServerhandleConsume(w.http.RespnsWriter, r *http.Request) {
    var req ConsumeRequest
    err := json.NewDecoder(r.body).Decode(&req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    record, err := s.log.Read(req.Offset)
    if err == ErrOffsetNotFound {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    res := ConsumeResponse{Record: record}
    err = json.NewEncoder(w).Encode(res)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}
        
