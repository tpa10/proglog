package server

import (
    "encoding/json"
    "net/http"

    "github.com/gorilla/mux"
)

// Used to create an instance of an HTTP server
func NewHTTPServer(addr string) *http.Server {
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

type ConsumeRequest struct {
    Offset uint64 `json:"offset"`
}

type ConsumeResponse struct {
    Record Record `json:"record"`
}

// Handle request to create a new entry.
func (s *httpServer) handleProduce( w http.ResponseWriter, r *http.Request) {
    var req ProduceRequest

    /*
     * Unmarshal the JSON defining the new data to be appended
     *  into a Record structure.
     */
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    /*
     * Append the new record to the end of the log, retrieving the offset
     *  (think "key") into the log file where the new record exists.
     */
    off, err := s.Log.Append(req.Record)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    /*
     * Format the offset (aka response data) into a struct
     */
    res := ProduceResponse{Offset: off}

    /*
     * Marshal the response data into JSON format directly into the 
     *  http response writer (w)
     */
    err = json.NewEncoder(w).Encode(res)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

// Handle request to fetch a record.
func (s *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
    var req ConsumeRequest

    /*
     * Un-Marshall the JSON formatted request data (offset/key) into a struct
     */
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    /*
     * Use the provided offset to try and read a record.
     */
    record, err := s.Log.Read(req.Offset)
    
    /*
     * If the offset doesn't exist, we have a "Page not found" condition
     *  User, not operational or software error, so we give them the old
     *  "404" finger and move on.
     */
    if err == ErrOffsetNotFound {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    /*
     * If we get a different error (e.g. read, permissions, etc) 
     *  then we have some kind of opertional problem or a bug,
     *  so we blow them a "500" snotgram and move on. 
     *  (Prepare for the support call)
     */
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Format the returned data into a struct
    res := ConsumeResponse{Record: record}

    // Marshal the struct into JSON and send it to the http response writer
    err = json.NewEncoder(w).Encode(res)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}
        
